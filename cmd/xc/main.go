package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/joe-davidson1802/xc/models"
	"github.com/joe-davidson1802/xc/parser"
	"github.com/joe-davidson1802/xc/run"
	"github.com/posener/complete"
)

var (
	version = ""
)

type config struct {
	version, help, short, md bool
	filename                 string
}

var cfg = config{}

func flags() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Println("xc - list tasks")
		fmt.Println("xc [task...] - run tasks")
		flag.PrintDefaults()
	}
	flag.BoolVar(&cfg.version, "version", false, "show xc version")
	flag.BoolVar(&cfg.help, "help", false, "shows xc usage")
	flag.BoolVar(&cfg.help, "h", false, "shows xc usage")
	flag.StringVar(&cfg.filename, "file", "README.md", "specify markdown file that contains tasks")
	flag.StringVar(&cfg.filename, "f", "README.md", "specify markdown file that contains tasks")
	flag.BoolVar(&cfg.short, "short", false, "list task names in a short format")
	flag.BoolVar(&cfg.short, "s", false, "list task names in a short format")
	flag.BoolVar(&cfg.md, "md", false, "print the markdown for a task rather than running it")
	flag.Parse()
}

func parse() (t models.Tasks, err error) {
	b, err := os.Open(cfg.filename)
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
	return
}

func printTask(t models.Task, maxLen int) {
	padLen := maxLen - len(t.Name)
	pad := strings.Repeat(" ", padLen)
	var desc []string
	for _, d := range t.Description {
		desc = append(desc, fmt.Sprintf("%s", d))
	}
	if len(t.DependsOn) > 0 {
		desc = append(desc, fmt.Sprintf("Requires:  %s", strings.Join(t.DependsOn, ", ")))
	}
	if len(desc) == 0 {
		desc = append(desc, t.Commands...)
	}
	fmt.Printf("    %s%s  %s\n", t.Name, pad, desc[0])
	for _, d := range desc[1:] {
		fmt.Printf("    %s  %s\n", strings.Repeat(" ", maxLen), d)
	}
}

func main() {
	flags()
	t, err := parse()
	if completion(t) {
		return
	}
	// xc -version
	if cfg.version {
		fmt.Printf("xc version: %s\n", getVersion())
		return
	}
	// xc -h / xc -help
	if cfg.help {
		flag.Usage()
		return
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tav := getArgs()
	// xc
	if len(tav) == 0 {
		printTasks(t)
		return
	}
	// xc -md task1 task2
	if cfg.md {
		if len(tav) == 0 {
			fmt.Println("md requires atleast 1 task")
			os.Exit(1)
		}
		for _, tv := range tav {
			ta, ok := t.Get(tv)
			if !ok {
				fmt.Printf("%s is not a task\n", tav[0])
			}
			ta.Display(os.Stdout)
		}
		return

	}
	// xc task1 task2
	for _, tav := range tav {
		err = t.ValidateDependencies(tav, []string{})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = run.RunTask(context.Background(), t, tav)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

}
func getArgs() []string {
	var (
		args = flag.Args()
	)
	return args
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
