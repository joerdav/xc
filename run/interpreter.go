package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var (
	shellShebangRe          = regexp.MustCompile(`^#!\s?/(usr/)?bin/(env\s+)?(sh|bash|mksh|bats|zsh)`)
	otherSupportedShebangRe = regexp.MustCompile(`^#!(.+)`)
)

type interpreter struct {
	shellRunner    func(context.Context, *interp.Runner, *syntax.File) error
	shebangRunner  func(*exec.Cmd) error
	tempFilePrefix string
}

func interpShellRunner(ctx context.Context, runner *interp.Runner, file *syntax.File) error {
	return runner.Run(ctx, file)
}

func cmdShebangRunner(cmd *exec.Cmd) error {
	return cmd.Run()
}

func newInterpreter() interpreter {
	return interpreter{
		shellRunner:    interpShellRunner,
		shebangRunner:  cmdShebangRunner,
		tempFilePrefix: "xc_",
	}
}

func (i interpreter) Execute(ctx context.Context, script string, env []string, args []string, dir, logPrefix string) error {
	interpreterCmd, interpreterArgs, text, ok := parseShebang(script)
	if !ok {
		return i.executeShell(ctx, script, env, args, dir, logPrefix)
	}
	return i.executeShebang(ctx, interpreterCmd, interpreterArgs, text, env, args, dir, logPrefix)
}

//nolint:gosec // accept that command is being executed here from outside of xc
func (i interpreter) executeShebang(
	ctx context.Context,
	interpreterCmd string,
	interpreterArgs []string,
	text string,
	env []string,
	args []string,
	dir string,
	logPrefix string,
) error {
	f, err := os.CreateTemp("", i.tempFilePrefix)
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
	stdin, stdout, stderr := stdFiles(logPrefix)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return i.shebangRunner(cmd)
}

func (i interpreter) executeShell(ctx context.Context, text string, env []string, args []string, dir, logPrefix string) error {
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
		interp.StdIO(stdFiles(logPrefix)),
		interp.Dir(dir),
		interp.Params(args...),
	)
	if err != nil {
		return fmt.Errorf("failed to compose script: %w", err)
	}
	return i.shellRunner(ctx, runner, file)
}

func parseShebang(script string) (interpreterCmd string, interpreterArgs []string, text string, ok bool) {
	if script == "" {
		return "", nil, "", false
	}
	lines := strings.Split(strings.TrimSpace(script), "\n")
	if shellShebangRe.MatchString(lines[0]) {
		return "", nil, "", false
	}
	if !otherSupportedShebangRe.MatchString(lines[0]) {
		return "", nil, "", false
	}
	shebang := lines[0]
	interpreter := strings.TrimPrefix(shebang, "#!")
	interpreterParts := strings.Fields(strings.TrimPrefix(interpreter, "/usr/bin/env "))
	interpreterCmd = interpreterParts[0]
	interpreterArgs = interpreterParts[1:]
	return interpreterCmd, interpreterArgs, strings.Join(lines[1:], "\n"), true
}

func stdFiles(prefix string) (io.Reader, io.Writer, io.Writer) {
	return os.Stdin, newPrefixLogger(os.Stdout, prefix), newPrefixLogger(os.Stderr, prefix)
}
