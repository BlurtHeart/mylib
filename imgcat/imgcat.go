package imgcat

import (
	"encoding/base64"
	"io"
	"strings"
)

//
func Copy(w io.Writer, r io.Reader) error {
	header := strings.NewReader("\033]1337;File=inline=1:")
	footer := strings.NewReader("\a\n")

	_, err := io.Copy(w, header)
	if err != nil {
		return err
	}

	wc := base64.NewEncoder(base64.StdEncoding, w)
	_, err = io.Copy(wc, r)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, footer)
	if err != nil {
		return err
	}
	return err
}

// NewWriter returns a new image writer
func NewWriter(w io.Writer) io.WriteCloser {
	pr, pw := io.Pipe()
	wc := &writer{pw: pw, done: make(chan struct{})}
	go func() {
		defer close(wc.done)
		Copy(w, pr)
	}()
	return wc
}

type writer struct {
	pw   *io.PipeWriter
	done chan struct{}
}

func (wc *writer) Write(data []byte) (int, error) {
	return wc.pw.Write(data)
}

func (wc *writer) Close() error {
	if err := wc.pw.Close(); err != nil {
		return err
	}
	<-wc.done
	return nil
}
