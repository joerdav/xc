package run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joerdav/xc/models"
)

const MAX_DEPS = 50

func runCmd(c *exec.Cmd) error {
	return c.Run()
}

type Runner struct {
	sep, cmdRunner, flag string
	runner               func(*exec.Cmd) error
	tasks                models.Tasks
}

func NewRunner(ts models.Tasks) (runner Runner, err error) {
	runner = Runner{
		sep:       ";",
		cmdRunner: "bash",
		flag:      "-c",
		runner:    runCmd,
		tasks:     ts,
	}
	if runtime.GOOS == "windows" {
		runner.sep = "&&"
		runner.cmdRunner = "cmd"
		runner.flag = "/C"
	}
	for _, t := range ts {
		err = runner.ValidateDependencies(t.Name, []string{})
		if err != nil {
			return
		}
	}
	return
}

func (r *Runner) Run(ctx context.Context, name string) error {
	task, ok := r.tasks.Get(name)
	if !ok {
		return fmt.Errorf("task %s not found", name)
	}
	for _, t := range task.DependsOn {
		err := r.Run(ctx, t)
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
	cmds := strings.Join(cmdl, r.sep)
	cmd := exec.Command(r.cmdRunner, r.flag, cmds)
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
	err = r.runner(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) ValidateDependencies(task string, prevTasks []string) error {
	if len(prevTasks) >= MAX_DEPS {
		return fmt.Errorf("max dependency depth of %d reached", MAX_DEPS)
	}
	// Check exists
	t, ok := r.tasks.Get(task)
	if !ok {
		return fmt.Errorf("task %s not found", task)
	}
	if t.ParsingError != "" {
		return fmt.Errorf("task %s has a parsing error: %s", task, t.ParsingError)
	}
	for _, t := range t.DependsOn {
		st, ok := r.tasks.Get(t)
		if !ok {
			return fmt.Errorf("task %s not found", t)
		}
		for _, pt := range prevTasks {
			if pt == st.Name {
				return fmt.Errorf("task %s contians a circular dependency", t)
			}
		}
		err := r.ValidateDependencies(st.Name, append([]string{st.Name}, prevTasks...))
		if err != nil {
			return err
		}
	}
	return nil
}
