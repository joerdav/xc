package models

import (
	_ "embed"
	"io"
	"os"
	"strings"
	"testing"
)

var gnTempl = Task{
	Name:   "generate-templ",
	Script: "go run -mod=mod github.com/a-h/templ/cmd/templ generate\ngo mod tidy\n",
}
var gnTranslations = Task{
	Name:   "generate-translations",
	Script: "go run ./i18n/generate\n",
}
var gnAll = Task{
	Name: "generate-all",
	DependsOn: []string{
		"generate-templ",
		"generate-translations",
	},
	DepsBehaviour: DependencyBehaviourAsync,
}

//go:embed testdata/generate-templ.md
var gnTemplExpected string

//go:embed testdata/generate-translations.md
var gnTranslationsExpected string

//go:embed testdata/generate-all.md
var gnAllExpected string

//go:embed testdata/generate-templ.json
var gnTemplExpectedJSON string

//go:embed testdata/generate-translations.json
var gnTranslationsExpectedJSON string

//go:embed testdata/generate-all.json
var gnAllExpectedJSON string

func assertTask(t *testing.T, expected, actual string) {
	t.Helper()
	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Fatalf("model want=%s got=%s", expected, actual)
	}
}

func taskToMarkdown(task Task) string {
	r, w, _ := os.Pipe()
	defer w.Close()
	defer r.Close()
	task.Display(w)
	w.Close()
	res, _ := io.ReadAll(r)
	return string(res)
}

func taskToJSON(task Task) string {
	r, w, _ := os.Pipe()
	defer w.Close()
	defer r.Close()
	task.DisplayJSON(w)
	w.Close()
	res, _ := io.ReadAll(r)
	return string(res)
}

func TestDisplay(t *testing.T) {
	assertTask(t, gnTemplExpected, taskToMarkdown(gnTempl))
	assertTask(t, gnTranslationsExpected, taskToMarkdown(gnTranslations))
	assertTask(t, gnAllExpected, taskToMarkdown(gnAll))
}

func TestDisplayJSON(t *testing.T) {
	assertTask(t, gnTemplExpectedJSON, taskToJSON(gnTempl))
	assertTask(t, gnTranslationsExpectedJSON, taskToJSON(gnTranslations))
	assertTask(t, gnAllExpectedJSON, taskToJSON(gnAll))
}
