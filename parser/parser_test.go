package parser

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joerdav/xc/models"
)

//go:embed testdata/example.md
var s string

//go:embed testdata/notasks.md
var e string

func TestParseFile(t *testing.T) {
	p, err := NewParser(strings.NewReader(s))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
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
	_, err := NewParser(strings.NewReader(e))
	if err.Error() != "no Tasks section found" {
		t.Fatalf("expected error %v got: %v", "no Tasks section found", err)
	}

}
