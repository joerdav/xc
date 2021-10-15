package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/joe-davidson1802/xc/parser"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

var (
	version = ""
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		pflag.PrintDefaults()
	}

	var (
		versionFlag  bool
		helpFlag     bool
		listFlag     bool
		listWideFlag bool
		fileName     string
	)

	pflag.BoolVar(&versionFlag, "version", false, "show xc version")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "shows xc usage")
	pflag.StringVarP(&fileName, "file", "f", "README.md", "specify markdown file that contains tasks")
	pflag.BoolVarP(&listFlag, "list", "l", false, "list tasks")
	pflag.BoolVar(&listWideFlag, "lw", false, "list tasks and commands")
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
	if listFlag || listWideFlag {
		table := tablewriter.NewWriter(os.Stdout)
		headers := []string{"task", "description"}
		if listWideFlag {
			headers = append(headers, "command")
		}
		table.SetHeader(headers)
		for _, ta := range t {
			row := []string{ta.Name, ta.Description}
			if listWideFlag {
				row = append(row, ta.Command)
			}
			table.Append(row)
		}
		table.Render()
		return
	}
	tav, _ := getArgs()
	fmt.Println("running tasks")
	for _, tav := range tav {
		fmt.Println(tav)
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
