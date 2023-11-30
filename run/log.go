package run

import (
	"bytes"
	"fmt"
	"io"
)

var newLine = byte('\n')

type prefixLogger struct {
	w      io.Writer
	buf    *bytes.Buffer
	prefix []byte
}

func newPrefixLogger(w io.Writer, prefix string) *prefixLogger {
	streamer := &prefixLogger{
		w:      w,
		buf:    bytes.NewBuffer([]byte("")),
		prefix: []byte(fmt.Sprintf("%s｜ ", prefix)),
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

	return l.out((p))
}

func (l *prefixLogger) outputLines() error {
	for {
		line, err := l.buf.ReadBytes(newLine)

		if len(line) > 0 {
			if bytes.HasSuffix(line, []byte{newLine}) {
				if err := l.out(line); err != nil {
					return err
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

func (l *prefixLogger) out(p []byte) error {
	if len(p) < 1 {
		return nil
	}

	_, err := l.w.Write(append(l.prefix, p...))
	return err
}
