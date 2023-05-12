package main

import (
	"flag"
	"github.com/go-logr/stdr"
	api "github.com/progimage/pkg/apiserver"
	"github.com/progimage/pkg/events"
	img "github.com/progimage/pkg/image"
	"github.com/progimage/pkg/transformations"
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

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-sigs
		// gracefully shutdown transformer service and http server when a signal is received
		cancelFunc()
	}()

	em := events.NewInMemoryEvents()
	localFS := img.LocalFS{
		BasePath: opts.LocalFSBasePath,
	}

	// Init and run Image Transformers
	transformImage := transformations.NewLocalTransformImage(ctx, logger, localFS, em)
	transformImage.Run()

	localIDGenerator := &img.IdGeneratorMemory{}
	localFileMgr := img.NewFileManager(localFS, em, localIDGenerator)
	svc := api.V1Service{
		Uploader:   localFileMgr,
		Downloader: localFileMgr,
		Logger:     logger,
	}

	// Init and run http server
	server := api.NewAPIServer(api.NewAPIServerOpts{
		BindPort: opts.Port,
		Done:     ctx,
	}, svc)
	server.Run()
}

func newStdLogger(flags int) stdr.StdLogger {
	return log.New(os.Stdout, "", flags)
}
