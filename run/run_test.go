package run

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joerdav/xc/models"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name               string
		runtime            string
		tasks              models.Tasks
		taskName           string
		expectedCmds       []string
		expectedRunError   bool
		expectedParseError bool
	}{
		{
			name:    "given an invalid task should not run command",
			runtime: "darwin",
			tasks: []models.Task{
				{
					Name: "mytask",
				},
			},
			taskName:         "fake",
			expectedCmds:     []string{},
			expectedRunError: true,
		},
		{
			name:    "given a valid command should run",
			runtime: "darwin",
			tasks: []models.Task{
				{
					Name: "mytask",
					Commands: []string{
						"somecmd",
					},
				},
			},
			taskName:     "mytask",
			expectedCmds: []string{`bash -c echo "somecmd";somecmd`},
		},
		{
			name:    "given a circular task should not run",
			runtime: "darwin",
			tasks: []models.Task{
				{
					Name:      "mytask",
					DependsOn: []string{"mytask2"},
				},
				{
					Name:      "mytask2",
					DependsOn: []string{"mytask"},
				},
			},
			taskName:           "mytask",
			expectedCmds:       []string{},
			expectedParseError: true,
		},
		{
			name:    "given a valid command with dep should run",
			runtime: "darwin",
			tasks: []models.Task{
				{
					Name: "mytask",
					Commands: []string{
						"somecmd",
					},
				},
				{
					Name: "mytask2",
					Commands: []string{
						"somecmd2",
					},
					DependsOn: []string{"mytask"},
				},
			},
			taskName:     "mytask2",
			expectedCmds: []string{`bash -c echo "somecmd";somecmd`, `bash -c echo "somecmd2";somecmd2`},
		},
		{
			name:    "given a valid command with dep should run on windows",
			runtime: "windows",
			tasks: []models.Task{
				{
					Name: "mytask",
					Commands: []string{
						"somecmd",
					},
				},
				{
					Name: "mytask2",
					Commands: []string{
						"somecmd2",
					},
					DependsOn: []string{"mytask"},
				},
			},
			taskName:     "mytask2",
			expectedCmds: []string{`cmd /C echo "somecmd"&&somecmd`, `cmd /C echo "somecmd2"&&somecmd2`},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cmds := []string{}
			runner, err := NewRunner(tt.tasks, tt.runtime)
			fmt.Println("results", err, tt.expectedParseError)
			if (err != nil) != tt.expectedParseError {
				t.Fatalf("expected error %v, got %v", tt.expectedParseError, err)
			}
			if err != nil {
				return
			}
			runner.runner = func(c *exec.Cmd) error {
				cmds = append(cmds, strings.Join(c.Args, " "))
				return nil
			}
			err = runner.Run(context.Background(), tt.taskName)
			if (err != nil) != tt.expectedRunError {
				t.Fatalf("expected error %v, got %v", tt.expectedRunError, err)
			}
			if diff := cmp.Diff(tt.expectedCmds, cmds); diff != "" {
				t.Fatalf("invalid commands received, %s", diff)
			}
		})
	}
}
