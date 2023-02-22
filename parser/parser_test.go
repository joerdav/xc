package parser

import (
	_ "embed"
	"errors"
	"strings"
	"testing"

	"github.com/joerdav/xc/models"
)

//go:embed testdata/example.md
var s string

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
	if strings.Join(expected.DependsOn, ",") != strings.Join(actual.DependsOn, ",") {
		t.Fatalf("requires want=%v got=%v", expected.DependsOn, actual.DependsOn)
	}
	if strings.Join(expected.Inputs, ",") != strings.Join(actual.Inputs, ",") {
		t.Fatalf("inputs want=%v got=%v", expected.Inputs, actual.Inputs)
	}
}

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

func TestParseFileNoTasks(t *testing.T) {
	_, err := NewParser(strings.NewReader(e))
	if !errors.Is(err, ErrNoTasksTitle) {
		t.Fatalf("expected error %v got: %v", "no Tasks section found", err)
	}
}

func TestMultipleDirs(t *testing.T) {
	p, _ := NewParser(strings.NewReader("dir: some dir"))
	p.currTask.Dir = "an existing dir"
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
`))
	_, err := p.parseTask()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestUnTerminatedCodeBlock(t *testing.T) {
	p, _ := NewParser(strings.NewReader(`
# Tasks
## a task
` + "```" + `
some code
`))
	_, err := p.parseTask()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestMultipleCodeBlocks(t *testing.T) {
	p, _ := NewParser(strings.NewReader("```\ncode\n```"))
	p.currTask.Script = "an existing script"
	err := p.parseCodeBlock()
	if err == nil {
		t.Fatal("expected error got nil")
	}
}

func TestParseAttribute(t *testing.T) {
	tests := []struct {
		name            string
		in              string
		expectNotOk     bool
		expectEnv       string
		expectDir       string
		expectDependsOn string
		expectInputs    string
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
			p, _ := NewParser(strings.NewReader(tt.in))
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
		})
	}
}
