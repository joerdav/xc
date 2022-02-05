package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/joe-davidson1802/xc/models"
)

// TRIM_VALUES are the characters that should be ignored in titles and attributes
const TRIM_VALUES = "_*` "

type parser struct {
	scanner    *bufio.Scanner
	tasks      models.Tasks
	currTask   models.Task
	tasksLevel int
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

func (p *parser) parseTitle(advance bool) (ok bool, level int, text string) {
	t := strings.TrimSpace(p.scanner.Text())
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
	p.scanner.Scan()
	return
}

type AttributeType int

const (
	AttributeTypeEnv AttributeType = iota
	AttributeTypeDir
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
	s := strings.Split(p.scanner.Text(), ":")
	if len(s) < 2 {
		return
	}
	ty, ok := attMap[strings.ToLower(strings.Trim(s[0], TRIM_VALUES))]
	if !ok {
		return
	}
	rest := strings.Join(s[1:], " ")
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
	p.scanner.Scan()
	return
}

func (p *parser) parseCodeBlock() (ok bool, err error) {
	t := p.scanner.Text()
	if len(t) < 3 || t[:3] != "```" {
		return
	}
	if len(p.currTask.Commands) > 0 {
		err = fmt.Errorf("command block already exists for task %s", p.currTask.Name)
		return
	}
	var ended bool
	for p.scanner.Scan() {
		if len(p.scanner.Text()) >= 3 && p.scanner.Text()[:3] == "```" {
			ended = true
			break
		}
		if strings.TrimSpace(p.scanner.Text()) != "" {
			p.currTask.Commands = append(p.currTask.Commands, p.scanner.Text())
		}
	}
	if !ended {
		err = fmt.Errorf("command block in task %s was not ended", p.currTask.Name)
		return
	}
	p.scanner.Scan()
	return
}

func (p *parser) parseTask() (ok bool, err error) {
	p.currTask = models.Task{}
	for {
		tok, level, text := p.parseTitle(true)
		if !tok || level > p.tasksLevel+1 {
			if !p.scanner.Scan() {
				break
			}
			continue
		}
		if level <= p.tasksLevel {
			return
		}
		p.currTask.Name = strings.Trim(text, TRIM_VALUES)
		break
	}
	if p.scanner.Err() != nil {
		err = p.scanner.Err()
		return
	}
	for {
		ok, err = p.parseAttribute()
		if ok {
			continue
		}
		if err != nil {
			return
		}
		ok, err = p.parseCodeBlock()
		if ok {
			continue
		}
		if err != nil {
			return
		}
		tok, level, _ := p.parseTitle(false)
		if tok && level <= p.tasksLevel {
			break
		}
		if tok && level == p.tasksLevel+1 {
			ok = true
			break
		}
		if strings.TrimSpace(p.scanner.Text()) != "" {
			p.currTask.Description = append(p.currTask.Description, strings.Trim(p.scanner.Text(), TRIM_VALUES))
		}
		if !p.scanner.Scan() {
			break
		}
	}
	if len(p.currTask.Commands) < 1 && len(p.currTask.DependsOn) < 1 {
		err = fmt.Errorf("task %s has no commands or required tasks %v", p.currTask.Name, p.currTask)
		return
	}
	p.tasks = append(p.tasks, p.currTask)
	return
}

func NewParser(r io.Reader) (p parser, err error) {
	p.scanner = bufio.NewScanner(r)
	for p.scanner.Scan() {
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
