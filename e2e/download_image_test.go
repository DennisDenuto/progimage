package e2e

import (
	"context"
	v1 "github.com/progimage/pkg/models/v1"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDownloadImage(t *testing.T) {
	// wait for http server to start
	require.Eventually(t, func() bool {
		return strings.Contains(string(httpServerSession.Buffer().Contents()), "GET    /api/v1/image/:imageID")
	}, 10*time.Second, 1*time.Second)

	b, w, err := createDogImageRequest(t)
	resp, err := client.UploadImageWithBody(context.Background(), w.FormDataContentType(), &b)
	require.NoError(t, err)

	response, err := v1.ParseUploadImageResponse(resp)
	require.NoError(t, err)

	uploadId := response.JSON200.Id

	t.Run("should be able to download uploaded dog image", func(t *testing.T) {
		image, err := client.DownloadImage(context.Background(), *uploadId)
		require.NoError(t, err)

		actualDogImage, err := io.ReadAll(image.Body)
		require.NoError(t, err)
		assetDogImage, err := os.ReadFile(dogImage)
		require.NoError(t, err)

		require.Equal(t, assetDogImage, actualDogImage)
	})

	t.Run("should be able to download compressed dog image", func(t *testing.T) {
		uploadId := *uploadId + ".gzip"
		image, err := client.DownloadImage(context.Background(), uploadId)
		require.NoError(t, err)

		actualDogImage, err := io.ReadAll(image.Body)
		require.NoError(t, err)
		assetDogImage, err := os.ReadFile(dogImageGzip)
		require.NoError(t, err)

		require.Equal(t, assetDogImage, actualDogImage)
	})
}
