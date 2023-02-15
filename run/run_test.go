package run

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/joerdav/xc/models"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name               string
		tasks              models.Tasks
		taskName           string
		err                error
		expectedRunError   bool
		expectedParseError bool
		expectedTasksRun   int
	}{
		{
			name: "given an invalid task should not run command",
			tasks: []models.Task{
				{
					Name: "mytask",
				},
			},
			taskName:         "fake",
			expectedRunError: true,
		},
		{
			name: "given a valid command should run",
			tasks: []models.Task{
				{
					Name:   "mytask",
					Script: "somecmd",
				},
			},
			taskName:         "mytask",
			expectedTasksRun: 1,
		},
		{
			name: "given a circular task should not run",
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
			expectedParseError: true,
		},
		{
			name: "given invalid script should fail",
			tasks: []models.Task{
				{
					Name:   "mytask",
					Script: "[[ ]]",
				},
			},
			taskName:         "mytask",
			expectedRunError: true,
		},
		{
			name: "given a valid command with deps only should run",
			tasks: []models.Task{
				{
					Name:   "mytask",
					Script: "somecmd",
				},
				{
					Name:      "mytask2",
					DependsOn: []string{"mytask"},
				},
			},
			taskName:         "mytask2",
			expectedTasksRun: 1,
		},
		{
			name: "given first task fails stop",
			tasks: []models.Task{
				{
					Name:   "mytask",
					Script: "somecmd",
				},
				{
					Name:      "mytask2",
					Script:    "somecmd2",
					Dir:       ".",
					DependsOn: []string{"mytask"},
				},
			},
			taskName:         "mytask2",
			err:              errors.New("some error"),
			expectedTasksRun: 1,
			expectedRunError: true,
		},
		{
			name: "given a valid command with dep should run",
			tasks: []models.Task{
				{
					Name:   "mytask",
					Script: "somecmd",
				},
				{
					Name:      "mytask2",
					Script:    "somecmd2",
					Dir:       ".",
					DependsOn: []string{"mytask"},
				},
			},
			taskName:         "mytask2",
			expectedTasksRun: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			runner, err := NewRunner(tt.tasks)
			fmt.Println("results", err, tt.expectedParseError)
			if (err != nil) != tt.expectedParseError {
				t.Fatalf("expected error %v, got %v", tt.expectedParseError, err)
			}
			if err != nil {
				return
			}
			runs := 0
			runner.scriptRunner = func(ctx context.Context, runner *interp.Runner, node syntax.Node) error {
				runs++
				return tt.err
			}
			err = runner.Run(context.Background(), tt.taskName, nil)
			if (err != nil) != tt.expectedRunError {
				t.Fatalf("expected error %v, got %v", tt.expectedRunError, err)
			}
			if runs != tt.expectedTasksRun {
				t.Fatalf("expected %d task runs got %d", tt.expectedTasksRun, runs)
			}
		})
	}
}
