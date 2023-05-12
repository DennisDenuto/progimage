package image

import "io"

type Downloader interface {
	Download(id string) (io.ReadCloser, error)
}
