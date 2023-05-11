package apiserver

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/progimage/image"
	v1 "github.com/progimage/models/v1"
	"io"
	"log"
	"net/http"
)

const streamChunkBytes = 32 * 1024

type NewAPIServerOpts struct {
	// BindPort is the port on which to serve HTTPS with authentication and authorization
	BindPort int

	APIRequestTimeout int

	stopCh chan struct{}
}

type APIServer struct {
	server *http.Server
	stopCh chan struct{}
	logger logr.Logger
	svc    v1Service
}

func (s *APIServer) Run() {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	v1.RegisterHandlersWithOptions(engine, s.svc, v1.GinServerOptions{
		BaseURL: "/api/v1",
	})

	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", 8080)); err != nil {
		panic(fmt.Sprintf("exited server unexpectedly %v", err))
	}
}

func NewAPIServer(opts NewAPIServerOpts) APIServer {
	return APIServer{
		svc: v1Service{
			uploader: image.NewLocalFileUploader("/tmp/"),
		},
	}
}

type v1Service struct {
	uploader image.Uploader
}

func (v v1Service) UploadImage(ctx *gin.Context) {
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
	uploadID, err := v.uploader.Upload(image.ImageRequest{
		Name:     file.Filename,
		Body:     f,
		Metadata: map[string]string{},
	})
	if err != nil {
		panic(err.Error())
	}

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", uploadID))
}

func (v v1Service) DownloadImage(ctx *gin.Context, imageID string) {
	// Stream the response in 32kB chunks
	reader := bytes.NewBufferString("")

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
