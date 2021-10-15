package models

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name        string
	Description string
	Command     string
}

type Tasks []Task

func (ts Tasks) Run(ctx context.Context, tsname string) error {
	var task *Task
	for _, t := range ts {
		if strings.ToLower(tsname) == strings.ToLower(t.Name) {
			task = &t
			break
		}
	}
	if task == nil {
		return fmt.Errorf("task %s not found", tsname)
	}
	cmd := exec.CommandContext(ctx, task.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
