package image

import (
	events2 "github.com/progimage/pkg/events"
	"io"
	"time"
)

type FS interface {
	Write(string, io.Reader) error
	Get(string) (io.ReadCloser, error)
}

type FileManager struct {
	fs FS

	uploadId     UploadId
	eventManager events2.FileEvents
}

func NewLocalFileManager(fs FS, em *events2.InMemoryEvents) *FileManager {
	return &FileManager{
		fs:           fs,
		uploadId:     &IdGeneratorMemory{},
		eventManager: em,
	}
}

func (l *FileManager) Download(id string) (io.ReadCloser, error) {
	return l.fs.Get(id)
}

func (l *FileManager) Upload(request ImageRequest) (ImageResponse, error) {
	id, err := l.uploadId.Generate()
	if err != nil {
		return ImageResponse{}, err
	}

	l.eventManager.Emit(id, events2.State{
		ID:        id,
		Val:       "UPLOAD_START",
		Timestamp: time.Now().UTC(),
	})

	go l.upload(id, request)

	return ImageResponse{UploadId: id}, nil
}

func (l *FileManager) upload(id string, request ImageRequest) {
	err := l.fs.Write(id, request.Body)
	if err != nil {
		l.eventManager.Emit(id, events2.State{
			ID:        id,
			Val:       "UPLOAD_ERROR",
			Timestamp: time.Now().UTC(),
			Metadata: map[string]string{
				"error": err.Error(),
			},
		})
		return
	}

	l.eventManager.Emit(id, events2.State{
		ID:        id,
		Val:       "UPLOAD_SUCCESS",
		Timestamp: time.Now().UTC(),
	})
}
