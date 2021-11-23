package models

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const MAX_DEPS = 50

var (
	COMMAND_SEP    string = ";"
	COMMAND_RUNNER        = "bash"
	COMMAND_FLAG          = "-c"
)

func init() {
	if runtime.GOOS == "windows" {
		COMMAND_SEP = "&&"
		COMMAND_RUNNER = "cmd"
		COMMAND_FLAG = "/C"
	}
}

type Task struct {
	Name         string
	Description  []string
	Commands     []string
	Dir          string
	Env          []string
	DependsOn    []string
	ParsingError string
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
	var cmdl []string
	for _, c := range task.Commands {
		if strings.TrimSpace(c) == "" {
			continue
		}
		cmdl = append(cmdl, fmt.Sprintf(`echo "%s"`, c), c)
	}
	if len(task.Commands) == 0 {
		return nil
	}
	cmds := strings.Join(cmdl, COMMAND_SEP)
	cmd := exec.Command(COMMAND_RUNNER, COMMAND_FLAG, cmds)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	for _, e := range task.Env {
		cmd.Env = append(cmd.Env, e)
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = path
	if task.Dir != "" {
		cmd.Dir = task.Dir
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
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
