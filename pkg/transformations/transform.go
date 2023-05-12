package transformations

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/progimage/pkg/events"
	image2 "github.com/progimage/pkg/image"
	"golang.org/x/net/context"
	"io"
)

type TransformImage struct {
	eventManager events.FileEvents
	fs           image2.FS
	transformers []Transformer
	logger       logr.Logger

	stop context.Context
}

func NewLocalTransformImage(ctx context.Context, logger logr.Logger, fs image2.LocalFS, em events.FileEvents) *TransformImage {
	return &TransformImage{
		eventManager: em,
		fs:           fs,
		transformers: []Transformer{
			CompressFile{
				Logger: logger,
			},
		},
		stop:   ctx,
		logger: logger,
	}
}

type Transformer interface {
	Transform([]byte, map[string]string) (io.Reader, string, error)
}

func (t TransformImage) Run() {
	go func() {
		for {
			select {
			case e := <-t.eventManager.WatchAll():
				if e.Val != "UPLOAD_SUCCESS" {
					continue
				}
				imageBytes, err := t.getImageContents(e)
				if err != nil {
					t.logger.Error(err, "Unable to get image contents", "id", e.ID)
					continue
				}

				for _, transformer := range t.transformers {
					transformedImage, suffix, err := transformer.Transform(imageBytes, e.Metadata)
					if err != nil {
						t.logger.Error(err, "unable to transform image.", e.ID)
						continue
					}

					err = t.fs.Write(fmt.Sprintf("%s.%s", e.ID, suffix), transformedImage)
					if err != nil {
						t.logger.Error(err, "unable to save transformed image.", e.ID)
						continue
					}
				}
			case <-t.stop.Done():
				t.logger.Info("Stopping transformations")
				return
			}
		}
	}()
}

func (t TransformImage) getImageContents(e events.State) ([]byte, error) {
	r, err := t.fs.Get(e.ID)
	if err != nil {
		return nil, err
	}
	imageBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}
