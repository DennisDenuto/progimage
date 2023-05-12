package image_test

import (
	"bytes"
	"github.com/progimage/pkg/events"
	image2 "github.com/progimage/pkg/image"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	baseDir := t.TempDir()
	fs := image2.LocalFS{baseDir}
	uploader := image2.NewFileManager(fs, events.NewInMemoryEvents(), &image2.IdGeneratorMemory{})

	uploadResponse, err := uploader.Upload(image2.ImageRequest{
		Name: "some-file",
		Body: bytes.NewBuffer([]byte{}),
		Metadata: map[string]string{
			"filetype": "png",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, uploadResponse.UploadId)

	t.Run("should eventually write file to local disk", func(t *testing.T) {
		require.Eventually(t, func() bool {
			_, err := os.Stat(filepath.Join(baseDir, uploadResponse.UploadId))
			return err == nil
		}, 30*time.Second, 1*time.Second)
	})
}
