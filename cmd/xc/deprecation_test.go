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
