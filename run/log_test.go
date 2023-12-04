//nolint:errcheck
package run

import (
	"bytes"
	"strings"
	"testing"
)

func setExpect(prefixB []byte, s string) string {
	prefix := string(prefixB)
	if len(prefix) == 0 {
		return s
	}
	if s == "" {
		return s
	}

	s = strings.ReplaceAll(s, "\n", "\n"+prefix)
	s = strings.TrimSuffix(s, prefix)

	return prefix + s
}

func TestPrefixLogger_Write(t *testing.T) {
	tests := map[string]string{"empty": "", "prefix": "prefix"}

	for name, prefix := range tests {
		t.Run(name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)

			l := newPrefixLogger(w, prefix)
			defer l.Close()

			l.Write([]byte("hello world"))

			// No new line so it should be empty
			expect := setExpect(l.prefix, "")
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}

			// Write a new line
			l.Write([]byte("\n"))
			expect = setExpect(l.prefix, "hello world\n")
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}

			// Write a line with a new line
			l.Write([]byte("foo bar\n"))
			expect = setExpect(l.prefix, "hello world\nfoo bar\n")
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}
		})
	}
}
