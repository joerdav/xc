package run

import "testing"

func TestIsShell(t *testing.T) {
	t.Run("no shebang assume shell", func(t *testing.T) {
		if !isShell("") {
			t.Fatal("expected true")
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
			if !isShell(she) {
				t.Errorf("%q expected true", she)
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
			if isShell(she) {
				t.Errorf("%q expected false", she)
			}
		}
	})
}
