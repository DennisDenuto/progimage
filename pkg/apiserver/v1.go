package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	img "github.com/progimage/pkg/image"
	v1 "github.com/progimage/pkg/models/v1"
	"io"
	"log"
	"net/http"
)

type V1Service struct {
	Uploader   img.Uploader
	Downloader img.Downloader
	Logger     logr.Logger
}

func (v V1Service) UploadImage(ctx *gin.Context) {
	// check content-length to prevent large files being uploaded

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		v.Logger.Error(err, "unable to get form file")
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		v.Logger.Error(err, "unable to open file")
		return
	}
	log.Println(file.Filename)
	// Upload the file to specific dst.
	uploadID, err := v.Uploader.Upload(img.ImageRequest{
		Name: file.Filename,
		Body: f,
		Metadata: map[string]string{
			"format": "image/png",
		},
	})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		v.Logger.Error(err, "unable to upload file to backend storage")
		return
	}

	ctx.JSON(http.StatusOK, v1.Artifact{
		Id: &uploadID.UploadId,
	})
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
