package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/joe-davidson1802/xc/parser"
	"github.com/spf13/pflag"
)

var (
	version = ""
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		fmt.Println("xc - list tasks")
		fmt.Println("xc [task...] - run tasks")
		pflag.PrintDefaults()
	}

	var (
		versionFlag bool
		helpFlag    bool
		fileName    string
	)

	pflag.BoolVar(&versionFlag, "version", false, "show xc version")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows xc usage")
	pflag.StringVarP(&fileName, "file", "f", "README.md", "specify markdown file that contains tasks")
	pflag.Parse()

	if versionFlag {
		fmt.Printf("xc version: %s\n", getVersion())
		return
	}

	if helpFlag {
		pflag.Usage()
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
	tav, _ := getArgs()
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
			fmt.Printf("    %s%s  %s\n", n.Name, pad, n.Description)
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
func getArgs() (tasksAndVars, cliArgs []string) {
	var (
		args          = pflag.Args()
		doubleDashPos = pflag.CommandLine.ArgsLenAtDash()
	)

	if doubleDashPos != -1 {
		tasksAndVars = args[:doubleDashPos]
		cliArgs = args[doubleDashPos:]
	} else {
		tasksAndVars = args
	}

	return
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
