package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/joe-davidson1802/xc/parser"
	"github.com/posener/complete"
)

var (
	version = ""
)

func completion(fileName string) bool {
	cmp := complete.New("xc", complete.Command{
		GlobalFlags: complete.Flags{
			"-version": complete.PredictNothing,
			"-h":       complete.PredictNothing,
			"-help":    complete.PredictNothing,
			"-f":       complete.PredictFiles("*.md"),
			"-file":    complete.PredictFiles("*.md"),
		},
	})
	b, err := os.ReadFile(fileName)
	if err == nil {
		f := string(b)
		t, err := parser.ParseFile(f)
		if err != nil {
			return false
		}
		s := make(map[string]complete.Command)
		for _, ta := range t {
			s[ta.Name] = complete.Command{}
		}
		cmp.Command.Sub = s
	}
	cmp.CLI.InstallName = "complete"
	cmp.CLI.UninstallName = "uncomplete"
	cmp.AddFlags(nil)

	flag.Parse()

	return cmp.Complete()
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Println("xc - list tasks")
		fmt.Println("xc [task...] - run tasks")
		flag.PrintDefaults()
	}

	var (
		versionFlag bool
		helpFlag    bool
		fileName    string
	)

	flag.BoolVar(&versionFlag, "version", false, "show xc version")
	flag.BoolVar(&helpFlag, "help", false, "shows xc usage")
	flag.BoolVar(&helpFlag, "h", false, "shows xc usage")
	flag.StringVar(&fileName, "file", "README.md", "specify markdown file that contains tasks")
	flag.StringVar(&fileName, "f", "README.md", "specify markdown file that contains tasks")

	if completion(fileName) {
		return
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	f := string(b)
	t, err := parser.ParseFile(f)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if versionFlag {
		fmt.Printf("xc version: %s\n", getVersion())
		return
	}

	if helpFlag {
		flag.Usage()
		return
	}
	tav := getArgs()
	if len(tav) == 0 {
		fmt.Println("tasks:")
		maxLen := 0
		for _, n := range t {
			if len(n.Name) > maxLen {
				maxLen = len(n.Name)
			}
		}
		for _, n := range t {
			padLen := maxLen - len(n.Name)
			pad := strings.Repeat(" ", padLen)
			var desc []string
			if n.ParsingError != "" {
				desc = append(desc, fmt.Sprintf("Parsing Error: %s", n.ParsingError))
			}
			for _, d := range n.Description {
				desc = append(desc, fmt.Sprintf("%s", d))
			}
			if len(n.DependsOn) > 0 {
				desc = append(desc, fmt.Sprintf("Requires:  %s", strings.Join(n.DependsOn, ", ")))
			}
			if len(desc) == 0 {
				desc = append(desc, n.Commands...)
			}
			fmt.Printf("    %s%s  %s\n", n.Name, pad, desc[0])
			for _, d := range desc[1:] {
				fmt.Printf("    %s  %s\n", strings.Repeat(" ", maxLen), d)
			}
		}
		return
	}
	for _, tav := range tav {
		err = t.ValidateDependencies(tav, []string{})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = t.Run(context.Background(), tav)
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
