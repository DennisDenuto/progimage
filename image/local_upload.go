package image

import (
	"github.com/progimage/events"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalFileUploader struct {
	basePath string

	uploadId     UploadId
	eventManager events.FileEvents
}

func NewLocalFileUploader(basePath string) *LocalFileUploader {
	return &LocalFileUploader{
		basePath: basePath,
		uploadId: &IdGeneratorMemory{},
		eventManager: &events.InMemoryEvents{
			Store: map[string][]events.State{},
		},
	}
}

func (l *LocalFileUploader) Upload(request ImageRequest) (ImageResponse, error) {
	id, err := l.uploadId.Generate()
	if err != nil {
		return ImageResponse{}, err
	}

	l.eventManager.Emit(id, events.State{
		Val:       "UPLOAD_START",
		Timestamp: time.Now().UTC(),
	})

	go l.upload(id, request)

	return ImageResponse{UploadId: id}, nil
}

func (l *LocalFileUploader) upload(id string, request ImageRequest) {
	create, err := os.Create(filepath.Join(l.basePath, id))
	if err != nil {
		l.eventManager.Emit(id, events.State{
			Val:       "UPLOAD_ERROR",
			Timestamp: time.Now().UTC(),
			Metadata: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}

	_, err = io.Copy(create, request.Body)
	if err != nil {
		l.eventManager.Emit(id, events.State{
			Val:       "UPLOAD_ERROR",
			Timestamp: time.Now().UTC(),
			Metadata: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}

	l.eventManager.Emit(id, events.State{
		Val:       "UPLOAD_SUCCESS",
		Timestamp: time.Now().UTC(),
	})
}
