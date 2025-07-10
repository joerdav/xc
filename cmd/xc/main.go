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
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/joerdav/xc/models"
	"github.com/joerdav/xc/parser/parsemd"
	"github.com/joerdav/xc/run"
	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/install"
	"github.com/posener/complete/v2/predict"
)

//go:embed usage.txt
var usage string

// ErrNoMarkdownFile will be returned if no markdown file is found in the cwd or any parent directories.
var ErrNoMarkdownFile = errors.New("no xc compatible documentation file found")

type config struct {
	version, help, short, display, noTTY, complete, uncomplete bool
	filename, heading                                          string
}

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
	flag.BoolVar(&cfg.version, "V", false, "show xc version")

	flag.BoolVar(&cfg.help, "help", false, "show xc usage")
	flag.BoolVar(&cfg.help, "h", false, "show xc usage")

	flag.StringVar(&cfg.heading, "heading", "", "specify the heading for xc tasks")
	flag.StringVar(&cfg.heading, "H", "", "specify the heading for xc tasks")

	flag.StringVar(&cfg.filename, "file", "", "specify a documentation file that contains tasks")
	flag.StringVar(&cfg.filename, "f", "", "specify a documentation file that contains tasks")

	flag.BoolVar(&cfg.short, "short", false, "list task names in a short format")
	flag.BoolVar(&cfg.short, "s", false, "list task names in a short format")

	flag.BoolVar(&cfg.display, "d", false, "print the plain text code of a task rather than running it")
	flag.BoolVar(&cfg.display, "display", false, "print the plain text code of a task rather than running it")

	flag.BoolVar(&cfg.complete, "complete", false, "install shell completion for xc")
	flag.BoolVar(&cfg.uncomplete, "uncomplete", false, "uninstall shell completion for xc")

	flag.BoolVar(&cfg.noTTY, "no-tty", false, "disable interactive picker")

	flag.Parse()
	return cfg
}

func parse(filename, heading string) (models.Tasks, string, error) {
	var specifiedHeading *string = nil
	if heading != "" {
		specifiedHeading = &heading
	}

	if filename != "" {
		return tryParse(filename, specifiedHeading)
	}
	curr, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, "", fmt.Errorf("error getting current directory: %w", err)
	}
	return searchUpForFile(curr, specifiedHeading)
}

func searchUpForFile(curr string, heading *string) (models.Tasks, string, error) {
	rm := filepath.Join(curr, "README.md")
	tasks, directory, err := tryParse(rm, heading)
	if err == nil {
		return tasks, directory, nil
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, parsemd.ErrNoTasksHeading) {
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
	return searchUpForFile(next, heading)
}

func tryParse(path string, heading *string) (models.Tasks, string, error) {
	directory := filepath.Dir(path)
	b, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("xc error opening file: %w", err)
	}
	p, err := parsemd.NewParser(b, heading)
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

func displayAndRunTasks(ctx context.Context, tasks models.Tasks, dir string, cfg config) error {
	if cfg.noTTY || cfg.short {
		printTasks(tasks, cfg.short)
		return nil
	}
	return interactivePicker(ctx, tasks, dir)
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
	ctx, cancel := context.WithCancel(context.Background())
	// handle SIGINT (control+c)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()
	cfg := flags()
	if cfg.uncomplete {
		return install.Uninstall("xc")
	}
	if cfg.complete {
		return install.Install("xc")
	}
	tasks, dir, err := parse(cfg.filename, cfg.heading)
	completion(tasks).Complete("xc")
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
		return displayAndRunTasks(ctx, tasks, dir, cfg)
	}
	ta, ok := tasks.Get(tav[0])
	if !ok {
		fmt.Printf("task \"%s\" not found\n", tav[0])
	}
	// xc -display task1
	if cfg.display {
		ta.Display(os.Stdout)
		return nil
	}
	// xc task1
	runner, err := run.NewRunner(tasks, dir)
	if err != nil {
		return fmt.Errorf("xc parse error: %w", err)
	}
	err = runner.Run(ctx, tav[0], tav[1:])
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

func completion(tasks models.Tasks) *complete.Command {
	return &complete.Command{
		Flags: map[string]complete.Predictor{
			"version": predict.Nothing,
			"V":       predict.Nothing,
			"h":       predict.Nothing,
			"help":    predict.Nothing,
			"f":       predict.Files("*.md"),
			"file":    predict.Files("*.md"),
			"s":       predict.Nothing,
			"short":   predict.Nothing,
			"d":       predict.Nothing,
			"display": predict.Nothing,
			"H":       predict.Nothing,
			"heading": predict.Nothing,
		},
		Sub: completeTasks(tasks),
	}
}

func completeTasks(tasks models.Tasks) map[string]*complete.Command {
	result := map[string]*complete.Command{}
	for _, t := range tasks {
		result[t.Name] = &complete.Command{
			Args: predict.Something,
		}
	}
	return result
}
