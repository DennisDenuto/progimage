package apiserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	image2 "github.com/progimage/pkg/image"
	"io"
	"log"
	"net/http"
)

type V1Service struct {
	Uploader   image2.Uploader
	Downloader image2.Downloader
}

func (v V1Service) UploadImage(ctx *gin.Context) {
	// check content-length to prevent large files being uploaded

	file, err := ctx.FormFile("file")
	if err != nil {
		panic(err.Error())
	}

	f, err := file.Open()
	if err != nil {
		panic(err.Error())
	}
	log.Println(file.Filename)
	// Upload the file to specific dst.
	uploadID, err := v.Uploader.Upload(image2.ImageRequest{
		Name: file.Filename,
		Body: f,
		Metadata: map[string]string{
			"format": "image/png",
		},
	})
	if err != nil {
		panic(err.Error())
	}

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", uploadID))
}

func (v V1Service) DownloadImage(ctx *gin.Context, imageID string) {
	// Stream the response in 32kB chunks
	reader, err := v.Downloader.Download(imageID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Stream(func(w io.Writer) bool {
		_, err := io.CopyN(w, reader, streamChunkBytes)
		if err != nil && err != io.EOF {
			ctx.Status(http.StatusInternalServerError)
			return false
		}
		return err != io.EOF
	})
}
