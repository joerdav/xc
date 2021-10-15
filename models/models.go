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
	fmt.Println(task.Name)
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
