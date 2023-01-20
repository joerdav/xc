package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/joerdav/xc/models"
)

// TRIM_VALUES are the characters that should be ignored in titles and attributes
const TRIM_VALUES = "_*` "

type parser struct {
	scanner               *bufio.Scanner
	tasks                 models.Tasks
	currTask              models.Task
	tasksLevel            int
	nextLine, currentLine string
	reachedEnd            bool
}

func (p *parser) Parse() (tasks models.Tasks, err error) {
	ok := true
	for ok {
		ok, err = p.parseTask()
		if err != nil || !ok {
			break
		}
	}
	tasks = p.tasks
	return
}

func (p *parser) scan() bool {
	if p.reachedEnd {
		return false
	}
	p.currentLine = p.nextLine
	if !p.scanner.Scan() {
		p.reachedEnd = true
		return true
	}
	p.nextLine = p.scanner.Text()
	return true
}

func (p *parser) parseAltTitle(advance bool) (ok bool, level int, text string) {
	t := strings.TrimSpace(p.currentLine)
	n := strings.TrimSpace(p.nextLine)
	if regexp.MustCompile("^-+$").MatchString(n) {
		ok = true
		level = 2
		text = t
	}
	if regexp.MustCompile("^=+$").MatchString(n) {
		ok = true
		level = 1
		text = t
	}
	if !advance || !ok {
		return
	}
	p.scan()
	p.scan()
	return
}

func (p *parser) parseTitle(advance bool) (ok bool, level int, text string) {
	ok, level, text = p.parseAltTitle(advance)
	if ok {
		return
	}
	t := strings.TrimSpace(p.currentLine)
	s := strings.Fields(t)
	if len(s) < 2 || len(s[0]) < 1 || strings.Count(s[0], "#") != len(s[0]) {
		return
	}
	ok = true
	level = len(s[0])
	text = strings.Join(s[1:], " ")
	if !advance {
		return
	}
	p.scan()
	return
}

// AttributeType represents metadata related to a Task.
//	# Tasks
//	## Task1
//	AttributeName: AttributeValue
//	```
//	script
//	```
type AttributeType int

const (
	// AttributeTypeEnv sets the environment variables for a Task.
	// It can be represented by an attribute with name `environment` or `env`.
	AttributeTypeEnv AttributeType = iota
	// AttributeTypeDir sets the working directory for a Task.
	// It can be represented by an attribute with name `directory` or `dir`.
	AttributeTypeDir
	// AttributeTypeReq sets the required Tasks for a Task, they will run
	// prior to the execution of the selected task.
	// It can be represented by an attribute with name `requires` or `req`.
	AttributeTypeReq
)

var attMap = map[string]AttributeType{
	"req":         AttributeTypeReq,
	"requires":    AttributeTypeReq,
	"env":         AttributeTypeEnv,
	"environment": AttributeTypeEnv,
	"dir":         AttributeTypeDir,
	"directory":   AttributeTypeDir,
}

func (p *parser) parseAttribute() (ok bool, err error) {
	a, rest, found := strings.Cut(p.currentLine, ":")
	if !found {
		return
	}
	ty, ok := attMap[strings.ToLower(strings.Trim(a, TRIM_VALUES))]
	if !ok {
		return
	}
	switch ty {
	case AttributeTypeReq:
		vs := strings.Split(rest, ",")
		for _, v := range vs {
			p.currTask.DependsOn = append(p.currTask.DependsOn, strings.Trim(v, TRIM_VALUES))
		}
	case AttributeTypeEnv:
		vs := strings.Split(rest, ",")
		for _, v := range vs {
			p.currTask.Env = append(p.currTask.Env, strings.Trim(v, TRIM_VALUES))
		}
	case AttributeTypeDir:
		if p.currTask.Dir != "" {
			err = fmt.Errorf("directory appears more than once for %s", p.currTask.Name)
			return
		}
		s := strings.Trim(rest, TRIM_VALUES)
		p.currTask.Dir = s
	}
	p.scan()
	return
}

func (p *parser) parseCodeBlock() (ok bool, err error) {
	t := p.currentLine
	if len(t) < 3 || t[:3] != "```" {
		return
	}
	if len(p.currTask.Commands) > 0 {
		err = fmt.Errorf("command block already exists for task %s", p.currTask.Name)
		return
	}
	var ended bool
	for p.scan() {
		if len(p.currentLine) >= 3 && p.currentLine[:3] == "```" {
			ended = true
			break
		}
		if strings.TrimSpace(p.currentLine) != "" {
			p.currTask.Commands = append(p.currTask.Commands, p.currentLine)
		}
	}
	if !ended {
		err = fmt.Errorf("command block in task %s was not ended", p.currTask.Name)
		return
	}
	p.scan()
	return
}

func (p *parser) findTaskTitle() (title string, done bool, err error) {
	for {
		tok, level, text := p.parseTitle(true)
		if !tok || level > p.tasksLevel+1 {
			if !p.scan() {
				return "", false, p.scanner.Err()
			}
			continue
		}
		if level <= p.tasksLevel {
			return "", true, nil
		}
		return strings.Trim(text, TRIM_VALUES), false, nil
	}
}

func (p *parser) parseTaskBody() (bool, error) {
	for {
		ok, err := p.parseAttribute()
		if ok {
			continue
		}
		if err != nil {
			return false, err
		}
		ok, err = p.parseCodeBlock()
		if ok {
			continue
		}
		if err != nil {
			return false, err
		}
		tok, level, _ := p.parseTitle(false)
		if tok && level <= p.tasksLevel {
			return false, nil
		}
		if tok && level == p.tasksLevel+1 {
			return true, nil
		}
		if strings.TrimSpace(p.currentLine) != "" {
			p.currTask.Description = append(p.currTask.Description, strings.Trim(p.currentLine, TRIM_VALUES))
		}
		if !p.scan() {
			return false, nil
		}
	}
}

func (p *parser) parseTask() (ok bool, err error) {
	p.currTask = models.Task{}
	title, done, err := p.findTaskTitle()
	if err != nil || done {
		return
	}
	p.currTask.Name = title
	ok, err = p.parseTaskBody()
	if err != nil {
		return
	}
	if len(p.currTask.Commands) < 1 && len(p.currTask.DependsOn) < 1 {
		err = fmt.Errorf("task %s has no commands or required tasks %v", p.currTask.Name, p.currTask)
		return
	}
	p.tasks = append(p.tasks, p.currTask)
	return
}

// NewParser will read from r until it finds a valid `tasks` block.
// If no block is found an error is returned.
func NewParser(r io.Reader) (p parser, err error) {
	p.scanner = bufio.NewScanner(r)
	for p.scan() {
		ok, level, text := p.parseTitle(true)
		if !ok || strings.ToLower(text) != "tasks" {
			continue
		}
		p.tasksLevel = level
		return
	}
	err = errors.New("no Tasks section found")
	return
}
