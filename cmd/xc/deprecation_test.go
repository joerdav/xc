package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/joerdav/xc/models"
)

func TestInteractiveWarning(t *testing.T) {
	tasks := models.Tasks{
		{
			Name:   "generate-templ",
			Script: "go run -mod=mod github.com/a-h/templ/cmd/templ generate\ngo mod tidy\n",
		},
		{
			Name:   "generate-translations",
			Script: "go run ./i18n/generate\n",
		},
		{
			Name: "generate-all",
			DependsOn: []string{
				"generate-templ",
				"generate-translations",
			},
			Interactive: true,
		},
	}
	orig := os.Stderr
	r, w, _ := os.Pipe()
	defer w.Close()
	defer r.Close()
	os.Stderr = w
	defer func() {
		os.Stderr = orig
	}()
	warnInteractive(tasks)
	w.Close()
	os.Stderr = orig
	res, _ := io.ReadAll(r)
	result := string(res)
	if !strings.Contains(result, "generate-all") {
		t.Fatal("warning does not contain task name")
	}
	if !strings.Contains(result, "Interactive") {
		t.Fatal("warning does not contain 'Interactive'")
	}
	if !strings.Contains(result, "https://github.com/joerdav/xc/") {
		t.Fatal("warning explanation does not contain repo link")
	}
}

func TestInteractiveWarningMultiple(t *testing.T) {
	tasks := models.Tasks{
		{
			Name:        "task1",
			Script:      "echo task1",
			Interactive: true,
		},
		{
			Name:   "task2",
			Script: "echo task2",
		},
		{
			Name:        "task3",
			Script:      "echo task3",
			Interactive: true,
		},
	}
	orig := os.Stderr
	r, w, _ := os.Pipe()
	defer w.Close()
	defer r.Close()
	os.Stderr = w
	defer func() {
		os.Stderr = orig
	}()
	warnInteractive(tasks)
	w.Close()
	os.Stderr = orig
	res, _ := io.ReadAll(r)
	result := string(res)

	// Should contain both task names in a single warning
	if !strings.Contains(result, "task1") {
		t.Fatal("warning does not contain task1")
	}
	if !strings.Contains(result, "task3") {
		t.Fatal("warning does not contain task3")
	}
	if !strings.Contains(result, "Interactive") {
		t.Fatal("warning does not contain 'Interactive'")
	}

	// Should only have one warning message (count occurrences of "warning:")
	warningCount := strings.Count(result, "warning: xc:")
	if warningCount != 1 {
		t.Fatalf("expected 1 warning message, got %d", warningCount)
	}
}
