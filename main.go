package main

import (
	"flag"
	"github.com/go-logr/stdr"
	"github.com/progimage/apiserver"
	"github.com/progimage/events"
	"github.com/progimage/image"
	"github.com/progimage/transformations"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	LocalFSBasePath string
	Port            int
}

func main() {
	logger := stdr.New(newStdLogger(log.Lshortfile))

	opts := Options{}
	flag.IntVar(&opts.Port, "port", 8080, "server bind port to use")
	flag.StringVar(&opts.LocalFSBasePath, "basePath", "/tmp/", "base path for the local fs to use.")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	em := events.NewInMemoryEvents()
	localFS := image.LocalFS{
		BasePath: opts.LocalFSBasePath,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-sigs
		// gracefully shutdown transformer service and http server when a signal is received
		cancelFunc()
	}()

	transformImage := transformations.NewLocalTransformImage(ctx, logger, localFS, em)
	transformImage.Run()

	localFileMgr := image.NewLocalFileManager(localFS, em)
	svc := apiserver.V1Service{
		Uploader:   localFileMgr,
		Downloader: localFileMgr,
	}

	server := apiserver.NewAPIServer(apiserver.NewAPIServerOpts{
		BindPort: opts.Port,
		Done:     ctx,
	}, svc)
	server.Run()
}

func newStdLogger(flags int) stdr.StdLogger {
	return log.New(os.Stdout, "", flags)
}
