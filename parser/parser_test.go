package parser

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/joerdav/xc/models"
)

//go:embed testdata/example.md
var s string

//go:embed testdata/marked-heading-next-line.md
var markedHeadingNextLine string

//go:embed testdata/marked-heading-same-line.md
var markedHeadingSameLine string

//go:embed testdata/till-eof.md
var tillEOF string

//go:embed testdata/notasks.md
var e string

func assertTask(t *testing.T, expected, actual models.Task) {
	t.Helper()
	if expected.Name != actual.Name {
		t.Fatalf("name want=%q got=%q", expected.Name, actual.Name)
	}
	if strings.Join(expected.Description, ",") != strings.Join(actual.Description, ",") {
		t.Fatalf("description want=%v got=%v", expected.Description, actual.Description)
	}
	if expected.Script != actual.Script {
		t.Fatalf("script want=%q got=%q", expected.Script, actual.Script)
	}
	if expected.Dir != actual.Dir {
		t.Fatalf("dir want=%q got=%q", expected.Dir, actual.Dir)
	}
	if expected.RequiredBehaviour != actual.RequiredBehaviour {
		t.Fatalf("Run want=%q got=%q", expected.RequiredBehaviour, actual.RequiredBehaviour)
	}
	if expected.DepsBehaviour != actual.DepsBehaviour {
		t.Fatalf("Run want=%q got=%q", expected.DepsBehaviour, actual.DepsBehaviour)
	}
	if strings.Join(expected.DependsOn, ",") != strings.Join(actual.DependsOn, ",") {
		t.Fatalf("requires want=%v got=%v", expected.DependsOn, actual.DependsOn)
	}
	if strings.Join(expected.Inputs, ",") != strings.Join(actual.Inputs, ",") {
		t.Fatalf("inputs want=%v got=%v", expected.Inputs, actual.Inputs)
	}
}

func TestParseFile(t *testing.T) {
	p, err := NewParser(strings.NewReader(s), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := models.Tasks{
		{Name: "list", Description: []string{"Lists files"}, Script: "ls\n"},
		{
			Name:        "list2",
			Description: []string{"Lists files"},
			Script:      "ls\n",
			Dir:         "./somefolder",
		},
		{
			Name:        "hello",
			Description: []string{"Print a message"},
			Script: `echo "Hello, world!"
echo "Hello, world2!"
`,
			Env:       []string{"somevar=val"},
			DependsOn: []string{"list", "list2"},
			Inputs:    []string{"FOO", "BAR"},
		},
		{
			Name:        "all-lists",
			Description: []string{"An example of a commandless task."},
			DependsOn:   []string{"list", "list2"},
		},
	}
	if len(result) != len(expected) {
		t.Fatalf("want %d tasks got %d", len(expected), len(result))
	}
	for i := range result {
		assertTask(t, expected[i], result[i])
	}
}

func TestParseFileToEOF(t *testing.T) {
	p, err := NewParser(strings.NewReader(tillEOF), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := models.Tasks{
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
			DepsBehaviour: models.DependencyBehaviourAsync,
		},
	}
	if len(result) != len(expected) {
		t.Fatalf("want %d tasks got %d", len(expected), len(result))
	}
	for i := range result {
		assertTask(t, expected[i], result[i])
	}
}

func TestParseFileNoTasks(t *testing.T) {
	_, err := NewParser(strings.NewReader(e), nil)
	if !errors.Is(err, ErrNoTasksHeading) {
		t.Fatalf("expected error %v got: %v", "no Tasks section found", err)
	}
}

func TestMultipleDirs(t *testing.T) {
	var p parser
	p.scanner = bufio.NewScanner(strings.NewReader("dir: some dir"))
	p.scan()
	p.scan()
	p.currTask.Dir = "an existing dir"
	_, err := p.parseAttribute()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestInvalidRun(t *testing.T) {
	var p parser
	p.scanner = bufio.NewScanner(strings.NewReader("run: never"))
	p.scan()
	p.scan()
	_, err := p.parseAttribute()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestCommandlessTask(t *testing.T) {
	p, _ := NewParser(strings.NewReader(`
# Tasks
## a task
## another task
`), nil)
	_, err := p.parseTask()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestRequiresOnlyTask(t *testing.T) {
	p, _ := NewParser(strings.NewReader(`
# Tasks
## a-task
requires: some-task
`), nil)
	_, err := p.parseTask()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHeadingCaseInsensitive(t *testing.T) {
	tests := []struct {
		mdHeading, parserHeading string
	}{
		{"Tasks", "Tasks"},
		{"Tasks", "tasks"},
		{"tasks", "Tasks"},
		{"tasks", "tasks"},
		{" Tasks", "Tasks"},
		{" Tasks", "tasks"},
		{" tasks", "Tasks"},
		{" tasks", "tasks"},
		{"Tasks", " Tasks"},
		{"Tasks", " tasks"},
		{"tasks", " Tasks"},
		{"tasks", " tasks"},
	}
	for _, tt := range tests {
		p, _ := NewParser(strings.NewReader(fmt.Sprintf(`
# %s
## a task
`+codeBlockStarter+`
some code
`+codeBlockStarter+`
`, tt.mdHeading)), &tt.parserHeading)
		_, err := p.parseTask()
		if err != nil {
			t.Fatal(err)
		}
		assertTask(t, models.Task{
			Name:   "a task",
			Script: "some code\n",
		}, p.currTask)
	}
}

func TestCustomHeadingByFlag(t *testing.T) {
	customHeading := "Out Of Scope"
	p, err := NewParser(strings.NewReader(s), &customHeading)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := models.Tasks{
		{
			Name:        "something",
			Description: []string{"Print a message"},
			Script:      "echo \"Hello, world!\"\n",
			DependsOn:   []string{"list", "list2"},
		},
	}
	if len(result) != len(expected) {
		t.Fatalf("want %d tasks got %d", len(expected), len(result))
	}
	for i := range result {
		assertTask(t, expected[i], result[i])
	}
}

func TestCustomHeadingByNextLineMarker(t *testing.T) {
	p, err := NewParser(strings.NewReader(markedHeadingNextLine), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := models.Tasks{
		{
			Name:        "hello",
			Description: []string{"Print a message"},
			Script:      "echo \"Hello, world!\"\n",
		},
	}
	if len(result) != len(expected) {
		t.Fatalf("want %d tasks got %d", len(expected), len(result))
	}
	for i := range result {
		assertTask(t, expected[i], result[i])
	}
}

func TestCustomHeadingBySameLineMarker(t *testing.T) {
	p, err := NewParser(strings.NewReader(markedHeadingSameLine), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := models.Tasks{
		{
			Name:        "hello",
			Description: []string{"Print a message"},
			Script:      "echo \"Hello, world!\"\n",
		},
	}
	if len(result) != len(expected) {
		t.Fatalf("want %d tasks got %d", len(expected), len(result))
	}
	for i := range result {
		assertTask(t, expected[i], result[i])
	}
}

func TestUnTerminatedCodeBlock(t *testing.T) {
	p, _ := NewParser(strings.NewReader(`
# Tasks
## a task
`+codeBlockStarter+`
some code
`), nil)
	_, err := p.parseTask()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestMultipleCodeBlocks(t *testing.T) {
	var p parser
	p.scanner = bufio.NewScanner(strings.NewReader("```\ncode\n```"))
	p.scan()
	p.scan()
	p.currTask.Script = "an existing script"
	err := p.parseCodeBlock()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestParseAttribute(t *testing.T) {
	tests := []struct {
		name                string
		in                  string
		expectNotOk         bool
		expectEnv           string
		expectDir           string
		expectDependsOn     string
		expectInputs        string
		expectBehaviour     models.RequiredBehaviour
		expectDepsBehaviour models.DepsBehaviour
	}{
		{
			name:      "given a basic Env, should parse",
			in:        "Env: my attribute",
			expectEnv: "my attribute",
		},
		{
			name:      "given environment attribute with mixed casing, should parse",
			in:        "EnvIronMent: my attribute",
			expectEnv: "my attribute",
		},
		{
			name:      "given Env with colons, should parse",
			in:        "Env: my:attribute",
			expectEnv: "my:attribute",
		},
		{
			name:      "given Env with formatting, should parse",
			in:        "Env: _*`my:attribute_*`",
			expectEnv: "my:attribute",
		},
		{
			name:            "given a basic req, should parse",
			in:              "req: my attribute",
			expectDependsOn: "my attribute",
		},
		{
			name:            "given requires attribute with mixed casing, should parse",
			in:              "ReqUiRES: my attribute",
			expectDependsOn: "my attribute",
		},
		{
			name:            "given req with colons, should parse",
			in:              "req: my:attribute",
			expectDependsOn: "my:attribute",
		},
		{
			name:            "given req with formatting, should parse",
			in:              "req: _*`my:attribute_*`",
			expectDependsOn: "my:attribute",
		},
		{
			name:         "given a basic Inputs, should parse",
			in:           "Inputs: my attribute",
			expectInputs: "my attribute",
		},
		{
			name:         "given inputs attribute with mixed casing, should parse",
			in:           "InpUts: my attribute",
			expectInputs: "my attribute",
		},
		{
			name:         "given Inputs with colons, should parse",
			in:           "Inputs: my:attribute",
			expectInputs: "my:attribute",
		},
		{
			name:         "given Inputs with formatting, should parse",
			in:           "Inputs: _*`my:attribute_*`",
			expectInputs: "my:attribute",
		},
		{
			name:      "given a basic dir, should parse",
			in:        "dir: my attribute",
			expectDir: "my attribute",
		},
		{
			name:      "given directory attribute with mixed casing, should parse",
			in:        "dIrECTORY: my attribute",
			expectDir: "my attribute",
		},
		{
			name:      "given dir with colons, should parse",
			in:        "dir: my:attribute",
			expectDir: "my:attribute",
		},
		{
			name:      "given dir with formatting, should parse",
			in:        "dir: _*`my:attribute_*`",
			expectDir: "my:attribute",
		},
		{
			name:            "given run always, should parse",
			in:              "run: always",
			expectBehaviour: models.RequiredBehaviourAlways,
		},
		{
			name:            "given run once, should parse",
			in:              "run: once",
			expectBehaviour: models.RequiredBehaviourOnce,
		},
		{
			name:            "given run once with formatting, should parse",
			in:              "run: _*`once`*_",
			expectBehaviour: models.RequiredBehaviourOnce,
		},
		{
			name:                "given runDeps sync, should parse",
			in:                  "runDeps: sync",
			expectDepsBehaviour: models.DependencyBehaviourSync,
		},
		{
			name:                "given runDeps async, should parse",
			in:                  "runDeps: async",
			expectDepsBehaviour: models.DependencyBehaviourAsync,
		},
		{
			name:                "given runDeps sync with formatting, should parse",
			in:                  "runDeps: _*`sync`*_",
			expectDepsBehaviour: models.DependencyBehaviourSync,
		},
		{
			name:        "given env with no colon, should not parse",
			in:          "env _*`my:attribute_*`",
			expectNotOk: true,
		},
		{
			name:        "given dir with no colon, should not parse",
			in:          "dir _*`my:attribute_*`",
			expectNotOk: true,
		},
		{
			name:        "given req with no colon, should not parse",
			in:          "req _*`my:attribute_*`",
			expectNotOk: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var p parser
			p.scanner = bufio.NewScanner(strings.NewReader(tt.in))
			p.scan()
			p.scan()
			ok, err := p.parseAttribute()
			if err != nil {
				t.Fatal(err)
			}
			if ok == tt.expectNotOk {
				t.Fatalf("ok=%v want=%v", ok, !tt.expectNotOk)
			}
			if tt.expectEnv != "" && p.currTask.Env[0] != tt.expectEnv {
				t.Fatalf("Env[0]=%s, want=%s", p.currTask.Env[0], tt.expectEnv)
			}
			if tt.expectDependsOn != "" && p.currTask.DependsOn[0] != tt.expectDependsOn {
				t.Fatalf("DependsOn[0]=%s, want=%s", p.currTask.DependsOn[0], tt.expectDependsOn)
			}
			if tt.expectInputs != "" && p.currTask.Inputs[0] != tt.expectInputs {
				t.Fatalf("Inputs[0]=%s, want=%s", p.currTask.Inputs[0], tt.expectInputs)
			}
			if tt.expectDir != "" && p.currTask.Dir != tt.expectDir {
				t.Fatalf("Dir=%s, want=%s", p.currTask.Dir, tt.expectDir)
			}
			if p.currTask.RequiredBehaviour != tt.expectBehaviour {
				t.Fatalf("got=%q, want=%q", p.currTask.RequiredBehaviour, tt.expectBehaviour)
			}
			if p.currTask.DepsBehaviour != tt.expectDepsBehaviour {
				t.Fatalf("got=%q, want=%q", p.currTask.DepsBehaviour, tt.expectDepsBehaviour)
			}
		})
	}
}

func BenchmarkParse10_000Tasks(b *testing.B) {
	var buf bytes.Buffer
	buf.WriteString(`
## Tasks
		`)
	for i := 0; i < 100; i++ {
		buf.WriteString(`
### task-` + fmt.Sprint(i) + `

Print a message

Requires: list, list2

` + "Env: `somevar=val`" + `
Inputs: FOO, BAR

` + codeBlockStarter + `
echo "Hello, world!"
echo "Hello, world2!"
` + codeBlockStarter)
	}
	file := buf.String()
	for i := 0; i < b.N; i++ {
		p, err := NewParser(strings.NewReader(file), nil)
		if err != nil {
			b.Fatal(err)
		}
		_, err = p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}
