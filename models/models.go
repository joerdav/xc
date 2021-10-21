package models

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MAX_DEPS = 50

type Task struct {
	Name        string
	Description string
	Command     string
	DependsOn   []string
}

type Tasks []Task

func (ts Tasks) Run(ctx context.Context, tsname string) error {
	task, ok := ts.Get(tsname)
	if !ok {
		return fmt.Errorf("task %s not found", tsname)
	}
	for _, t := range task.DependsOn {
		err := ts.Run(ctx, t)
		if err != nil {
			return err
		}
	}
	fmt.Println(task.Command)
	parts := strings.Split(task.Command, " ")
	cmd := exec.Command(parts[0])
	cmd.Args = parts[0:]
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = path
	return cmd.Run()
}

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
		return fmt.Errorf("task %s not found", t)
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
