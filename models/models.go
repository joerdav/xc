package models

import (
	"fmt"
	"io"
	"strings"
)

const MAX_DEPS = 50

type Task struct {
	Name         string
	Description  []string
	Commands     []string
	Dir          string
	Env          []string
	DependsOn    []string
	ParsingError string
}

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
	if len(t.Commands) > 0 {
		fmt.Fprintln(w, "```")
		for _, d := range t.Commands {
			fmt.Fprintln(w, d)
		}
		fmt.Fprintln(w, "```")
	}
}

type Tasks []Task

func (ts Tasks) Get(tsname string) (task Task, ok bool) {
	for _, t := range ts {
		if strings.ToLower(tsname) == strings.ToLower(t.Name) {
			ok = true
			task = t
			break
		}
	}
	return
}

func (ts Tasks) ValidateDependencies(task string, prevTasks []string) error {
	if len(prevTasks) >= MAX_DEPS {
		return fmt.Errorf("max dependency depth of %d reached", MAX_DEPS)
	}
	// Check exists
	t, ok := ts.Get(task)
	if !ok {
		return fmt.Errorf("task %s not found", task)
	}
	if t.ParsingError != "" {
		return fmt.Errorf("task %s has a parsing error: %s", task, t.ParsingError)
	}
	for _, t := range t.DependsOn {
		st, ok := ts.Get(t)
		if !ok {
			return fmt.Errorf("task %s not found", t)
		}
		for _, pt := range prevTasks {
			if pt == st.Name {
				return fmt.Errorf("task %s contians a circular dependency", t)
			}
		}
		err := ts.ValidateDependencies(st.Name, append([]string{st.Name}, prevTasks...))
		if err != nil {
			return err
		}
	}
	return nil
}
