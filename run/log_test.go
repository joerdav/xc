//nolint:errcheck
package run

import (
	"bytes"
	"testing"
)

func TestPrefixLogger_Write(t *testing.T) {
	tests := map[string]string{"empty": "", "prefix": "prefix"}

	for name, prefix := range tests {
		t.Run(name, func(t *testing.T) {
			w := bytes.NewBuffer(nil)

			l := newPrefixLogger(w, prefix)
			defer l.Close()

			l.Write([]byte("hello world"))

			// No new line so it should be empty
			expect := ""
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}

			// Write a new line
			l.Write([]byte("\n"))
			expect += string(l.prefix) + "hello world\n"
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}

			// Write a line with a new line
			l.Write([]byte("foo bar\n"))
			// Each line is prefixed
			expect += string(l.prefix) + "foo bar\n"
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}

			// Write a line with a new line with a red color
			red := "\033[31m"
			reset := "\033[m"
			l.Write([]byte(
				red + "this line is red\n" +
					"this line is still red" + reset + "\n" +
					"this line is reset\n",
			))
			// Each line specifies the output color
			expect += string(l.prefix) + red + "this line is red\n"
			expect += string(l.prefix) + red + "this line is still red" + reset + "\n"
			expect += string(l.prefix) + reset + "this line is reset\n"
			if w.String() != expect {
				t.Errorf("got %v, want %v", w.String(), expect)
			}
		})
	}
}
