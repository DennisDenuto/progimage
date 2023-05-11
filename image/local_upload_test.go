package image_test

import (
	"bytes"
	"github.com/progimage/image"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	baseDir := t.TempDir()
	uploader := image.NewLocalFileUploader(baseDir)

	uploadResponse, err := uploader.Upload(image.ImageRequest{
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
