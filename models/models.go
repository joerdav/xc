package models

import (
	"fmt"
	"io"
	"strings"
)

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
		if strings.EqualFold(tsname, t.Name) {
			ok = true
			task = t
			break
		}
	}
	return
}
