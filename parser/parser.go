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
var MissingCommand = errors.New("missing command or Requires")

var taskR = regexp.MustCompile(`(?i)^#+ +Tasks`)
var heading = regexp.MustCompile(`(?i)^#+`)
var commandDef = regexp.MustCompile(`(?i)^.+:.*$`)
var commandTitle = regexp.MustCompile(`(?i)^.+: *`)
var cleanName = regexp.MustCompile("(?i)[_*:` #]")
var codeBlock = regexp.MustCompile("(?i)^```.*$")
var deps = regexp.MustCompile("(?i)^Requires:.*$")
var dir = regexp.MustCompile("(?i)^Directory:.*$")
var env = regexp.MustCompile("(?i)^Env:.*$")
var trimValues = "_*` "

func isTask(text string, taskDepth int) bool {
	isHeading := heading.MatchString(text)
	comesUnderTasks := strings.Count(text, "#") == taskDepth+1
	return isHeading && comesUnderTasks
}
func isTaskSection(text string) bool {
	return taskR.MatchString(text)
}

func parseListValue(text string) (ss []string) {
	idx := strings.Index(text, ":") + 1
	s := strings.Trim(text[idx:], trimValues)
	ss = strings.Split(s, ",")
	for i := range ss {
		ss[i] = strings.Trim(ss[i], trimValues)
	}
	return
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
				if currentTask.Command == "" && len(currentTask.DependsOn) == 0 {
					currentTask.ParsingError = fmt.Sprintf("%v", MissingCommand)
				}
				ts = append(ts, *currentTask)
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
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			currentTask.Command += text
			continue
		}
		if env.MatchString(text) {
			ss := parseListValue(text)
			currentTask.Env = append(currentTask.Env, ss...)
			continue
		}
		if deps.MatchString(text) {
			ss := parseListValue(text)
			currentTask.DependsOn = append(currentTask.DependsOn, ss...)
			continue
		}
		if dir.MatchString(text) {
			if currentTask.Dir != "" {
				err = fmt.Errorf("directory appears more than once for %s", currentTask.Name)
				return
			}
			idx := strings.Index(text, ":") + 1
			s := text[idx:]
			s = strings.Trim(s, trimValues)
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
	if currentTask != nil {
		if currentTask.Command == "" && len(currentTask.DependsOn) == 0 {
			currentTask.ParsingError = fmt.Sprintf("%v", MissingCommand)
		}
		ts = append(ts, *currentTask)
	}
	return
}
