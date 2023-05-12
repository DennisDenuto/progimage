package image

import (
	"io"
	"os"
	"path/filepath"
)

type LocalFS struct {
	BasePath string
}

func (l LocalFS) Write(id string, reader io.Reader) error {
	create, err := os.Create(filepath.Join(l.BasePath, id))
	if err != nil {
		return err
	}

	_, err = io.Copy(create, reader)
	if err != nil {
		return err
	}

	return nil
}

func (l LocalFS) Get(id string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.BasePath, id))
}
