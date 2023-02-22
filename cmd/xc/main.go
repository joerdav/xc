package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/joerdav/xc/models"
	"github.com/joerdav/xc/parser"
	"github.com/joerdav/xc/run"
	"github.com/posener/complete"
)

// ErrNoMarkdownFile will be returned if no markdown file is found in the cwd or any parent directories.
var ErrNoMarkdownFile = errors.New("no xc compatible markdown file found")

type config struct {
	version, help, short, md bool
	filename                 string
}

//go:embed usage.txt
var usage string

var version = ""

func main() {
	if err := runMain(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func flags() config {
	var cfg config

	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.BoolVar(&cfg.version, "version", false, "show xc version")
	flag.BoolVar(&cfg.help, "help", false, "shows xc usage")
	flag.BoolVar(&cfg.help, "h", false, "shows xc usage")
	flag.StringVar(&cfg.filename, "file", "", "specify markdown file that contains tasks")
	flag.StringVar(&cfg.filename, "f", "", "specify markdown file that contains tasks")
	flag.BoolVar(&cfg.short, "short", false, "list task names in a short format")
	flag.BoolVar(&cfg.short, "s", false, "list task names in a short format")
	flag.BoolVar(&cfg.md, "md", false, "print the markdown for a task rather than running it")
	flag.Parse()
	return cfg
}

func parse(filename string) (models.Tasks, string, error) {
	if filename != "" {
		return tryParse(filename)
	}
	curr, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, "", fmt.Errorf("error getting current directory: %w", err)
	}
	return searchUpForFile(curr)
}

func searchUpForFile(curr string) (models.Tasks, string, error) {
	rm := filepath.Join(curr, "README.md")
	tasks, directory, err := tryParse(rm)
	if err == nil {
		return tasks, directory, nil
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, parser.ErrNoTasksTitle) {
		return nil, "", err
	}
	git := filepath.Join(curr, ".git")
	_, err = os.Stat(git)
	if err == nil {
		return nil, "", ErrNoMarkdownFile
	}
	next := filepath.Dir(curr)
	if strings.HasSuffix(next, string([]rune{filepath.Separator})) {
		return nil, "", ErrNoMarkdownFile
	}
	return searchUpForFile(next)
}

func tryParse(path string) (models.Tasks, string, error) {
	directory := filepath.Dir(path)
	b, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("xc error opening file: %w", err)
	}
	p, err := parser.NewParser(b)
	if err != nil {
		return nil, "", fmt.Errorf("xc parse error: %w", err)
	}
	tasks, err := p.Parse()
	if err != nil {
		return nil, "", fmt.Errorf("xc parse error: %w", err)
	}
	return tasks, directory, nil
}

func printTasks(tasks models.Tasks, short bool) {
	print := printTask
	if short {
		print = func(t models.Task, maxLen int) { fmt.Println(t.Name) }
	}
	maxLen := 0
	for _, n := range tasks {
		if len(n.Name) > maxLen {
			maxLen = len(n.Name)
		}
	}
	for _, n := range tasks {
		print(n, maxLen)
	}
}

func printTask(task models.Task, maxLen int) {
	padLen := maxLen - len(task.Name)
	pad := strings.Repeat(" ", padLen)
	desc := task.Description
	if len(task.DependsOn) > 0 {
		desc = append(desc, fmt.Sprintf("Requires:  %s", strings.Join(task.DependsOn, ", ")))
	}
	if len(desc) == 0 {
		desc = strings.Split(task.Script, "\n")
	}
	fmt.Printf("    %s%s  %s\n", task.Name, pad, desc[0])
	for _, d := range desc[1:] {
		fmt.Printf("    %s  %s\n", strings.Repeat(" ", maxLen), d)
	}
}

func runMain() error {
	cfg := flags()
	tasks, dir, err := parse(cfg.filename)
	if completion(tasks) {
		return nil
	}
	// xc -version
	if cfg.version {
		fmt.Printf("xc version: %s\n", getVersion())
		return nil
	}
	// xc -h / xc -help
	if cfg.help {
		flag.Usage()
		return nil
	}
	if err != nil {
		return err
	}
	tav := flag.Args()
	// xc
	if len(tav) == 0 {
		printTasks(tasks, cfg.short)
		return nil
	}
	ta, ok := tasks.Get(tav[0])
	if !ok {
		fmt.Printf("task \"%s\" not found\n", tav[0])
	}
	// xc -md task1
	if cfg.md {
		ta.Display(os.Stdout)
		return nil
	}
	// xc task1
	runner, err := run.NewRunner(tasks, dir)
	if err != nil {
		return fmt.Errorf("xc parse error: %w", err)
	}
	err = runner.Run(context.Background(), tav[0], tav[1:])
	if err != nil {
		return fmt.Errorf("xc: %w", err)
	}
	return nil
}

func getVersion() string {
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "unknown"
	}

	version = info.Main.Version
	if info.Main.Sum != "" {
		version += fmt.Sprintf(" (%s)", info.Main.Sum)
	}

	return version
}

func completion(tasks models.Tasks) bool {
	cmp := complete.New("xc", complete.Command{
		GlobalFlags: complete.Flags{
			"-version": complete.PredictNothing,
			"-h":       complete.PredictNothing,
			"-short":   complete.PredictNothing,
			"-help":    complete.PredictNothing,
			"-f":       complete.PredictFiles("*.md"),
			"-file":    complete.PredictFiles("*.md"),
		},
		Sub:   nil,
		Flags: nil,
		Args:  nil,
	})
	commands := make(map[string]complete.Command)
	for _, ta := range tasks {
		commands[ta.Name] = complete.Command{
			Sub:         nil,
			Flags:       nil,
			GlobalFlags: nil,
			Args:        nil,
		}
	}
	cmp.Command.Sub = commands
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)
	return cmp.Complete()
}
