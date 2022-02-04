package run

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joe-davidson1802/xc/models"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name          string
		tasks         models.Tasks
		taskName      string
		expectedCmds  []string
		expectedError bool
	}{
		{
			name: "given an invalid task should not run command",
			tasks: []models.Task{
				{
					Name: "mytask",
				},
			},
			taskName:      "fake",
			expectedCmds:  []string{},
			expectedError: true,
		},
		{
			name: "given a valid command should run",
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
			name: "given a valid command with dep should run",
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cmds := []string{}
			runner := NewRunner(tt.tasks)
			runner.runner = func(c *exec.Cmd) error {
				cmds = append(cmds, strings.Join(c.Args, " "))
				return nil
			}
			err := runner.Run(context.Background(), tt.taskName)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
			if diff := cmp.Diff(tt.expectedCmds, cmds); diff != "" {
				t.Errorf("invalid commands received, %s", diff)
			}
		})
	}
}
