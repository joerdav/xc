package parser

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joe-davidson1802/xc/models"
)

//go:embed testdata/example.md
var s string

//go:embed testdata/notasks.md
var e string

func TestParseFile(t *testing.T) {
	result, err := ParseFile(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected %d tasks, got %d", 2, len(result))
	}
	task1 := models.Task{
		Name:        "list",
		Description: []string{"Lists files"},
		Command:     "ls",
	}
	if diff := cmp.Diff(task1, result[0]); diff != "" {
		t.Errorf("task1 does not match expected: %s", diff)
	}
	task3 := models.Task{
		Name:        "hello",
		Description: []string{"Print a message"},
		Command:     `echo "Hello, world!"`,
		DependsOn:   []string{"list", "list2"},
	}
	if diff := cmp.Diff(task3, result[2]); diff != "" {
		t.Errorf("task2 does not match expected: %s", diff)
	}

}
func TestParseFileNoTasks(t *testing.T) {
	_, err := ParseFile(e)
	if err != NoTasksError {
		t.Fatalf("expected error %v got: %v", NoTasksError, err)
	}

}
