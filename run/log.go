package run

import (
	"bytes"
	"io"
	"regexp"
)

var (
	prefixColor = []byte("\033[0m")
	delimiter   = []byte("｜ ")
	newLine     = byte('\n')
	colorRegexp = regexp.MustCompile(`\033\[[0-9;]*m`)
)

type prefixLogger struct {
	w            io.Writer
	buf          *bytes.Buffer
	prefix       []byte
	currentColor []byte
}

func newPrefixLogger(w io.Writer, prefix string) *prefixLogger {
	streamer := &prefixLogger{
		w:   w,
		buf: bytes.NewBuffer([]byte("")),
	}
	if prefix != "" {
		p := make([]byte, 0, len(prefixColor)+len(prefix)+len(delimiter))
		p = append(p, prefixColor...)
		p = append(p, []byte(prefix)...)
		p = append(p, delimiter...)

		streamer.prefix = p
	}

	return streamer
}

func (l *prefixLogger) Write(p []byte) (n int, err error) {
	if n, err = l.buf.Write(p); err != nil {
		return
	}

	err = l.outputLines()
	return
}

func (l *prefixLogger) Close() error {
	if err := l.Flush(); err != nil {
		return err
	}
	l.buf = bytes.NewBuffer([]byte(""))
	return nil
}

func (l *prefixLogger) Flush() error {
	p := make([]byte, l.buf.Len())
	if _, err := l.buf.Read(p); err != nil {
		return err
	}

	return l.out((p), []byte{})
}

func (l *prefixLogger) outputLines() error {
	for {
		line, err := l.buf.ReadBytes(newLine)

		if len(line) > 0 {
			if bytes.HasSuffix(line, []byte{newLine}) {
				if err := l.out(line, l.currentColor); err != nil {
					return err
				}

				colors := colorRegexp.FindAll(line, -1)
				if len(colors) > 0 {
					l.currentColor = colors[len(colors)-1]
				}
			} else {
				// put back into buffer, it's not a complete line yet
				//  Close() or Flush() have to be used to flush out
				//  the last remaining line if it does not end with a newline
				if _, err := l.buf.Write(line); err != nil {
					return err
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *prefixLogger) out(p []byte, currentColor []byte) error {
	if len(p) < 1 {
		return nil
	}

	s := make([]byte, 0, len(l.prefix)+len(currentColor)+len(p))
	s = append(s, l.prefix...)
	s = append(s, currentColor...)
	s = append(s, p...)

	_, err := l.w.Write(s)
	return err
}
