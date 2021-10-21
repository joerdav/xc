package parser

import (
	"bufio"
	"errors"
	"regexp"
	"strings"

	"github.com/joe-davidson1802/xc/models"
)

var NoTasksError = errors.New("no tasks found")
var MissingCommand = errors.New("missing command")

var taskR = regexp.MustCompile(`^#+ +Tasks`)
var heading = regexp.MustCompile(`^#+`)
var commandDef = regexp.MustCompile(`^.+:.*$`)
var commandTitle = regexp.MustCompile(`^.+: *`)
var cleanName = regexp.MustCompile(`[_*: ]`)
var codeBlock = regexp.MustCompile("^```.*$")
var deps = regexp.MustCompile("^>.*$")

func ParseFile(f string) (ts models.Tasks, err error) {
	var foundTasksSection bool
	var level int
	var inCodeBlock bool
	var currentTask *models.Task
	scanner := bufio.NewScanner(strings.NewReader(f))
	for scanner.Scan() {
		if taskR.MatchString(scanner.Text()) {
			foundTasksSection = true
			level = strings.Count(scanner.Text(), "#")
			continue
		}
		if !foundTasksSection {
			continue
		}
		if heading.MatchString(scanner.Text()) && strings.Count(scanner.Text(), "#") <= level {
			break
		}
		if commandDef.MatchString(scanner.Text()) {
			if currentTask != nil {
				err = MissingCommand
				return
			}
			name := commandTitle.FindStringSubmatch(scanner.Text())[0]
			name = cleanName.ReplaceAllString(name, "")
			description := commandTitle.ReplaceAllString(scanner.Text(), "")
			currentTask = &models.Task{
				Name:        name,
				Description: description,
			}
			continue
		}
		if currentTask == nil {
			continue
		}
		if codeBlock.MatchString(scanner.Text()) {
			if inCodeBlock {
				ts = append(ts, *currentTask)
				currentTask = nil
			}
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			currentTask.Command += scanner.Text()
			continue
		}
		if deps.MatchString(scanner.Text()) {
			s := strings.ReplaceAll(scanner.Text(), ">", "")
			ss := strings.Split(s, ",")
			for i := range ss {
				ss[i] = strings.Trim(ss[i], " ")
			}
			currentTask.DependsOn = append(currentTask.DependsOn, ss...)
			continue
		}
	}
	if !foundTasksSection {
		err = NoTasksError
		return
	}
	return
}
