package parser

import (
	"bufio"
	"errors"
	"fmt"
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
var cleanName = regexp.MustCompile(`[_*: #]`)
var codeBlock = regexp.MustCompile("^```.*$")
var deps = regexp.MustCompile("^Requires:.*$")
var dir = regexp.MustCompile("^Directory:.*$")

func isTask(text string, taskDepth int) bool {
	isHeading := heading.MatchString(text)
	comesUnderTasks := strings.Count(text, "#") > taskDepth
	return isHeading && comesUnderTasks
}
func isTaskSection(text string) bool {
	return taskR.MatchString(text)
}

func ParseFile(f string) (ts models.Tasks, err error) {
	var foundTasksSection bool
	var taskLevel int
	var inCodeBlock bool
	var currentTask *models.Task
	scanner := bufio.NewScanner(strings.NewReader(f))
	for scanner.Scan() {
		text := scanner.Text()
		if isTaskSection(text) {
			foundTasksSection = true
			taskLevel = strings.Count(text, "#")
			continue
		}
		if !foundTasksSection {
			continue
		}
		if heading.MatchString(text) && strings.Count(text, "#") <= taskLevel {
			break
		}
		if isTask(text, taskLevel) {
			if currentTask != nil {
				err = fmt.Errorf("%v: near %s", MissingCommand, text)
				return
			}
			name := cleanName.ReplaceAllString(text, "")
			currentTask = &models.Task{
				Name: name,
			}
			continue
		}
		if currentTask == nil {
			continue
		}
		if codeBlock.MatchString(text) {
			if inCodeBlock {
				ts = append(ts, *currentTask)
				currentTask = nil
			}
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			currentTask.Command += text
			continue
		}
		if deps.MatchString(text) {
			s := strings.ReplaceAll(scanner.Text(), "Requires:", "")
			ss := strings.Split(s, ",")
			for i := range ss {
				ss[i] = strings.Trim(ss[i], " ")
			}
			currentTask.DependsOn = append(currentTask.DependsOn, ss...)
			continue
		}
		if dir.MatchString(text) {
			if currentTask.Dir != "" {
				err = fmt.Errorf("directory appears more than once for %s", currentTask.Name)
				return
			}
			s := strings.ReplaceAll(scanner.Text(), "Directory:", "")
			s = strings.Trim(s, " ")
			currentTask.Dir = s
			continue
		}
		if text != "" {
			currentTask.Description = append(currentTask.Description, text)
			continue
		}
	}
	if !foundTasksSection {
		err = NoTasksError
		return
	}
	return
}
