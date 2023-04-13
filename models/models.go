package models

import (
	"fmt"
	"io"
	"strings"
)

// Task represents a parsed Task.
type Task struct {
	Name              string
	Description       []string
	Script            string
	Dir               string
	Env               []string
	DependsOn         []string
	Inputs            []string
	ParsingError      string
	RequiredBehaviour RequiredBehaviour
}

// Display writes a Task as Markdown.
func (t Task) Display(w io.Writer) {
	fmt.Fprintf(w, "## %s\n\n", t.Name)
	for _, d := range t.Description {
		fmt.Fprintln(w, d)
		fmt.Fprintln(w)
	}
	if len(t.DependsOn) > 0 {
		fmt.Fprintln(w, "Requires:", strings.Join(t.DependsOn, ", "))
		fmt.Fprintln(w)
	}
	if t.Dir != "" {
		fmt.Fprintln(w, "Directory:", t.Dir)
		fmt.Fprintln(w)
	}
	if len(t.Env) > 0 {
		fmt.Fprintln(w, "Env:", strings.Join(t.Env, ", "))
		fmt.Fprintln(w)
	}
	if len(t.Inputs) > 0 {
		fmt.Fprintln(w, "Inputs:", strings.Join(t.Inputs, ", "))
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "Run:", t.RequiredBehaviour)
	fmt.Fprintln(w)
	if len(t.Script) > 0 {
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w, t.Script)
		fmt.Fprintln(w, "```")
	}
}

// Tasks is an alias type for []Task
type Tasks []Task

// Get returns a task by name, case insensitively.
func (ts Tasks) Get(tsname string) (task Task, ok bool) {
	for _, t := range ts {
		if strings.EqualFold(tsname, t.Name) {
			ok = true
			task = t
			break
		}
	}
	return
}

// RequiredBehaviour represents a tasks behaviour when
// required by another task.
// The default is RequiredBehaviourAlways
type RequiredBehaviour int

const (
	// RequiredBehaviourAlways should be used if the task is to be run every time it is required.
	RequiredBehaviourAlways RequiredBehaviour = iota
	// RequiredBehaviourOnce should be used if a task should be run once, even if required multiple times.
	RequiredBehaviourOnce
)

func (b RequiredBehaviour) String() string {
	if b == RequiredBehaviourOnce {
		return "once"
	}
	return "always"
}

func ParseRequiredBehaviour(s string) (RequiredBehaviour, bool) {
	switch strings.ToLower(s) {
	case "once":
		return RequiredBehaviourOnce, true
	case "always":
		return RequiredBehaviourAlways, true
	default:
		return 0, false
	}
}
