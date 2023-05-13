package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	img "github.com/progimage/pkg/image"
	v1 "github.com/progimage/pkg/models/v1"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type V1Service struct {
	Uploader   img.Uploader
	Downloader img.Downloader
	Logger     logr.Logger
}

func (v V1Service) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		v.Logger.Error(err, "unable to get form file")
		return
	}

	// basic validation. only allow certain image types.
	contentType := v.tryContentType(file)
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "application/octet-stream":
		v.Logger.Info("uploading file", "contentType", contentType)
	default:
		ctx.Status(http.StatusBadRequest)
		v.Logger.Error(err, "invalid image type uploaded", "contentType", contentType)
		return
	}

	f, err := file.Open()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		v.Logger.Error(err, "unable to open file")
		return
	}
	// Upload the file to specific dst.
	uploadID, err := v.Uploader.Upload(img.ImageRequest{
		Name: file.Filename,
		Body: f,
		Metadata: map[string]string{
			"timestamp":  time.Now().UTC().String(),
			"image_type": contentType,
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
	reader, err := v.Downloader.Download(imageID)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Stream(func(w io.Writer) bool {
		// Stream the response in 32kB chunks
		_, err := io.CopyN(w, reader, streamChunkBytes)
		if err != nil && err != io.EOF {
			ctx.Status(http.StatusInternalServerError)
			return false
		}
		return err != io.EOF
	})
}

func (v V1Service) tryContentType(file *multipart.FileHeader) string {
	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	open, err := file.Open()
	if err != nil {
		return "application/octet-stream"
	}
	_, err = open.Read(buff)
	if err != nil {
		return "application/octet-stream"
	}
	return http.DetectContentType(buff)
}
