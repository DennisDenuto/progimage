package transformations

import (
	"bytes"
	"compress/gzip"
	"github.com/go-logr/logr"
	"io"
)

type CompressFile struct {
	Logger logr.Logger
}

func (c CompressFile) Transform(s []byte, metadata map[string]string) (io.Reader, string, error) {
	c.Logger.Info("transforming: compressing image", "metadata", metadata)
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)
	defer func() {
		_ = gz.Close()
	}()

	if _, err := gz.Write(s); err != nil {
		return nil, "", err
	}

	return &b, "gzip", nil
}
