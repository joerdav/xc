package run

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var (
	shellShebangRe          = regexp.MustCompile(`^#!\s?/(usr/)?bin/(env\s+)?(sh|bash|mksh|bats|zsh)\s`)
	otherSupportedShebangRe = regexp.MustCompile(`^#!(.+)`)
)

type interpreter struct{}

func (interpreter) Execute(ctx context.Context, script string, env []string, args []string, dir string) error {
	if isShell(script) {
		return executeShell(ctx, script, env, args, dir)
	}
	return executeShebang(ctx, script, env, args, dir)
}

func executeShebang(ctx context.Context, text string, env []string, args []string, dir string) error {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	shebang := lines[0]
	interpreter := strings.TrimPrefix(shebang, "#!")
	interpreterParts := strings.Fields(strings.TrimPrefix(interpreter, "/usr/bin/env "))
	interpreterCmd := interpreterParts[0]
	interpreterArgs := interpreterParts[1:]
	text = strings.Join(lines[1:], "\n")
	d, err := ioutil.TempDir("", "xc_")
	if err != nil {
		return fmt.Errorf("failed to create execution dir")
	}
	defer os.Remove(d)
	f, err := ioutil.TempFile(d, "xc_")
	if err != nil {
		return fmt.Errorf("failed to create execution file")
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(text); err != nil {
		return fmt.Errorf("failed to write execution file")
	}
	interpreterArgs = append(interpreterArgs, f.Name())
	cmd := exec.CommandContext(ctx, interpreterCmd, append(interpreterArgs, args...)...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeShell(ctx context.Context, text string, env []string, args []string, dir string) error {
	if shellShebangRe.MatchString(text) {
		text = strings.Join(strings.Split(text, "\n")[1:], "\n")
	}
	var script bytes.Buffer
	if _, err := script.Write([]byte(scriptHeader)); err != nil {
		return fmt.Errorf("failed to write script header: %w", err)
	}
	if _, err := script.Write([]byte(text)); err != nil {
		return fmt.Errorf("failed to write script: %w", err)
	}
	file, err := syntax.NewParser().Parse(&script, "")
	if err != nil {
		return fmt.Errorf("failed to parse task: %w", err)
	}
	runner, err := interp.New(
		interp.Env(expand.ListEnviron(env...)),
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Dir(dir),
		interp.Params(args...),
	)
	if err != nil {
		return fmt.Errorf("failed to compose script: %w", err)
	}
	return runner.Run(ctx, file)
}

func isShell(script string) bool {
	lines := strings.Split(strings.TrimSpace(script), "\n")
	if len(lines) == 0 {
		return true
	}
	if shellShebangRe.MatchString(lines[0]) {
		return true
	}
	if !otherSupportedShebangRe.MatchString(lines[0]) {
		return true
	}
	return false
}
