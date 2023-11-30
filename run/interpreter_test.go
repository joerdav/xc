package run

import (
	"context"
	"os/exec"
	"testing"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type testInterpreter struct {
	interpreter
	shellRunnerCalled   bool
	shebangRunnerCalled bool
}

func newTestInterpreter() *testInterpreter {
	inter := &testInterpreter{}
	inter.interpreter = interpreter{
		shellRunner: func(context.Context, *interp.Runner, *syntax.File) error {
			inter.shellRunnerCalled = true
			return nil
		},
		shebangRunner: func(*exec.Cmd) error {
			inter.shebangRunnerCalled = true
			return nil
		},
	}
	return inter
}

func TestIsShell(t *testing.T) {
	t.Run("empty assume shell", func(t *testing.T) {
		ti := newTestInterpreter()
		if err := ti.Execute(context.Background(), "", nil, nil, "", ""); err != nil {
			t.Fatal(err)
		}
		if !ti.shellRunnerCalled {
			t.Fatal("expected shell call")
		}
		if ti.shebangRunnerCalled {
			t.Fatal("expected no shebang")
		}
	})
	t.Run("no shebang assume shell", func(t *testing.T) {
		ti := newTestInterpreter()
		if err := ti.Execute(context.Background(), "echo", nil, nil, "", ""); err != nil {
			t.Fatal(err)
		}
		if !ti.shellRunnerCalled {
			t.Fatal("expected shell call")
		}
		if ti.shebangRunnerCalled {
			t.Fatal("expected no shebang")
		}
	})
	t.Run("shell shebang should result in shell", func(t *testing.T) {
		shells := []string{
			"sh",
			"bash",
			"mksh",
			"bats",
			"zsh",
		}
		for _, s := range shells {
			she := "#!/usr/bin/env " + s + " "
			ti := newTestInterpreter()
			if err := ti.Execute(context.Background(), she, nil, nil, "", ""); err != nil {
				t.Fatal(err)
			}
			if !ti.shellRunnerCalled {
				t.Fatal("expected shell call")
			}
			if ti.shebangRunnerCalled {
				t.Fatal("expected no shebang")
			}
		}
	})
	t.Run("other shebang should not result in shell", func(t *testing.T) {
		shells := []string{
			"python",
			"node",
		}
		for _, s := range shells {
			she := "#!/usr/bin/env " + s + " "
			ti := newTestInterpreter()
			if err := ti.Execute(context.Background(), she, nil, nil, "", ""); err != nil {
				t.Fatal(err)
			}
			if ti.shellRunnerCalled {
				t.Fatal("expected no shell call")
			}
			if !ti.shebangRunnerCalled {
				t.Fatal("expected shebang")
			}
		}
	})
	t.Run("shell shebang with invalid bash script should fail", func(t *testing.T) {
		she := `#!/usr/bin/env bash

		func main() {
			print("hang on this isn't shell")
		}`
		ti := newTestInterpreter()
		if err := ti.Execute(context.Background(), she, nil, nil, "", ""); err == nil {
			t.Fatal("expected an error")
		}
		if ti.shellRunnerCalled {
			t.Fatal("expected no shell call")
		}
		if ti.shebangRunnerCalled {
			t.Fatal("expected no shebang")
		}
	})
	t.Run("error creating file should not execute", func(t *testing.T) {
		she := "#!/usr/bin/env python "
		ti := newTestInterpreter()
		ti.tempFilePrefix = "invalid/prefix"
		if err := ti.Execute(context.Background(), she, nil, nil, "", ""); err == nil {
			t.Fatal("expected an error")
		}
		if ti.shellRunnerCalled {
			t.Fatal("expected no shell call")
		}
		if ti.shebangRunnerCalled {
			t.Fatal("expected no shebang")
		}
	})
}
