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

var (
	version = ""
	cfg     = config{}
)

func main() {
	if err := runMain(); err != nil {
		fmt.Println("xc error: ", err.Error())
		os.Exit(1)
	}
}

func flags() {
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
}

func parse() (t models.Tasks, directory string, err error) {
	if cfg.filename != "" {
		return tryParse(cfg.filename)
	}
	curr, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, "", err
	}
	return searchUpForFile(curr)
}

func searchUpForFile(curr string) (t models.Tasks, directory string, err error) {
	rm := filepath.Join(curr, "README.md")
	t, directory, err = tryParse(rm)
	if err == nil {
		return t, directory, err
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

func tryParse(path string) (t models.Tasks, directory string, err error) {
	directory = filepath.Dir(path)
	if err != nil {
		return
	}
	b, err := os.Open(path)
	if err != nil {
		return
	}
	p, err := parser.NewParser(b)
	if err != nil {
		return
	}
	t, err = p.Parse()
	return
}

func printTasks(t models.Tasks) {
	print := printTask
	if cfg.short {
		print = func(t models.Task, maxLen int) { fmt.Println(t.Name) }
	}
	maxLen := 0
	for _, n := range t {
		if len(n.Name) > maxLen {
			maxLen = len(n.Name)
		}
	}
	for _, n := range t {
		print(n, maxLen)
	}
}

func printTask(t models.Task, maxLen int) {
	padLen := maxLen - len(t.Name)
	pad := strings.Repeat(" ", padLen)
	var desc []string
	desc = append(desc, t.Description...)
	if len(t.DependsOn) > 0 {
		desc = append(desc, fmt.Sprintf("Requires:  %s", strings.Join(t.DependsOn, ", ")))
	}
	if len(desc) == 0 {
		desc = strings.Split(t.Script, "\n")
	}
	fmt.Printf("    %s%s  %s\n", t.Name, pad, desc[0])
	for _, d := range desc[1:] {
		fmt.Printf("    %s  %s\n", strings.Repeat(" ", maxLen), d)
	}
}

func runMain() error {
	flags()
	t, dir, err := parse()
	if completion(t) {
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
		printTasks(t)
		return nil
	}
	// xc -md task1
	if cfg.md {
		if len(tav) == 0 {
			fmt.Println("a task is required")
			flag.Usage()
			os.Exit(1)
		}
		ta, ok := t.Get(tav[0])
		if !ok {
			fmt.Printf("%s is not a task\n", tav[0])
		}
		ta.Display(os.Stdout)
		return nil

	}
	// xc task1
	runner, err := run.NewRunner(t, dir)
	if err != nil {
		return err
	}
	err = runner.Run(context.Background(), tav[0], tav[1:])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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

func completion(t models.Tasks) bool {
	cmp := complete.New("xc", complete.Command{
		GlobalFlags: complete.Flags{
			"-version": complete.PredictNothing,
			"-h":       complete.PredictNothing,
			"-short":   complete.PredictNothing,
			"-help":    complete.PredictNothing,
			"-f":       complete.PredictFiles("*.md"),
			"-file":    complete.PredictFiles("*.md"),
		},
	})
	s := make(map[string]complete.Command)
	for _, ta := range t {
		s[ta.Name] = complete.Command{}
	}
	cmp.Command.Sub = s
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)
	return cmp.Complete()
}
