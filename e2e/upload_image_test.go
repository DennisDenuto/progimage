package e2e

import (
	"bytes"
	"context"
	v1 "github.com/progimage/pkg/models/v1"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"testing"
	"time"
)

const dogImage = "./assets/dog.jpeg"
const dogImageGzip = "./assets/dog.jpeg.gz"

func TestUploadImage(t *testing.T) {
	// wait for http server to start
	require.Eventually(t, func() bool {
		return strings.Contains(string(httpServerSession.Buffer().Contents()), "GET    /api/v1/image/:imageID")
	}, 10*time.Second, 1*time.Second)

	b, w, err := createDogImageRequest(t)

	resp, err := client.UploadImageWithBody(context.Background(), w.FormDataContentType(), &b)
	require.NoError(t, err)

	response, err := v1.ParseUploadImageResponse(resp)
	require.NotNil(t, response)
	require.NotNil(t, response.JSON200)

	require.GreaterOrEqual(t, *response.JSON200.Id, "0")
}

func createDogImageRequest(t *testing.T) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	dogImage, err := os.Open(dogImage)
	require.NoError(t, err)

	fw, err := w.CreateFormFile("file", dogImage.Name())
	require.NoError(t, err)

	_, err = io.Copy(fw, dogImage)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	return b, w, err
}
