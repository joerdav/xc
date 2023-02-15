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

func TestParseAttribute(t *testing.T) {
	tests := []struct {
		name            string
		in              string
		expectNotOk     bool
		expectEnv       string
		expectDir       string
		expectDependsOn string
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
				t.Error(err)
			}
			if ok == tt.expectNotOk {
				t.Errorf("ok=%v want=%v", ok, !tt.expectNotOk)
			}
			if tt.expectEnv != "" && p.currTask.Env[0] != tt.expectEnv {
				t.Errorf("Env[0]=%s, want=%s", p.currTask.Env[0], tt.expectEnv)
			}
			if tt.expectDependsOn != "" && p.currTask.DependsOn[0] != tt.expectDependsOn {
				t.Errorf("DependsOn[0]=%s, want=%s", p.currTask.DependsOn[0], tt.expectDependsOn)
			}
			if tt.expectDir != "" && p.currTask.Dir != tt.expectDir {
				t.Errorf("Dir=%s, want=%s", p.currTask.Dir, tt.expectDir)
			}
		})
	}
}
