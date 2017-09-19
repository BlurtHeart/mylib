package imgcat

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestNewWriter(t *testing.T) {
	filename := "gopher.png"
	if err := cat(filename); err != nil {
		fmt.Fprintf(os.Stderr, "could not cat %s, err:%v\n", filename, err)
	}
}

func cat(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return errors.New("could not open file")
	}
	defer fp.Close()

	wc := NewWriter(os.Stdout)
	if _, err := io.Copy(wc, fp); err != nil {
		return err
	}
	return wc.Close()
}
