package image

import (
	"io"
)

type ImageRequest struct {
	Name     string
	Body     io.Reader
	Metadata map[string]string
}

type ImageResponse struct {
	UploadId string
}

type Uploader interface {
	Upload(ImageRequest) (ImageResponse, error)
}

type UploadId interface {
	Generate() (string, error)
}
