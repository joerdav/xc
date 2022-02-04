package run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joe-davidson1802/xc/models"
)

func runCmd(c *exec.Cmd) error {
	return c.Run()
}

type Runner struct {
	sep, cmdRunner, flag string
	runner               func(*exec.Cmd) error
	tasks                models.Tasks
}

func NewRunner(ts models.Tasks) Runner {
	r := Runner{
		sep:       ";",
		cmdRunner: "bash",
		flag:      "-c",
		runner:    runCmd,
		tasks:     ts,
	}
	if runtime.GOOS == "windows" {
		r.sep = "&&"
		r.cmdRunner = "cmd"
		r.flag = "/C"
	}
	return r
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
