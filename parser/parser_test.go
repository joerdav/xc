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
	expected := models.Tasks{
		{Name: "list", Description: []string{"Lists files"}, Commands: []string{"ls"}},
		{
			Name:        "list2",
			Description: []string{"Lists files"},
			Commands:    []string{"ls"},
			Dir:         "./somefolder",
		},
		{
			Name:         "empty-task",
			Description:  []string{"Description but no task info"},
			ParsingError: "missing command or Requires",
		},
		{
			Name:        "hello",
			Description: []string{"Print a message"},
			Commands: []string{
				`echo "Hello, world!"`,
				`echo "Hello, world2!"`,
			},
			Env:       []string{"somevar=val"},
			DependsOn: []string{"list", "list2"},
		},
		{
			Name:        "all-lists",
			Description: []string{"An example of a commandless task."},
			DependsOn:   []string{"list", "list2"},
		},
	}
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("tasks does not match expected: %s", diff)
	}

}
func TestParseFileNoTasks(t *testing.T) {
	_, err := ParseFile(e)
	if err != NoTasksError {
		t.Fatalf("expected error %v got: %v", NoTasksError, err)
	}

}
