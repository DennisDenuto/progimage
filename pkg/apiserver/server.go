package apiserver

import (
	"context"
	"fmt"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/progimage/pkg/models/v1"
	"net/http"
	"os"
)

const streamChunkBytes = 32 * 1024

type NewAPIServerOpts struct {
	// BindPort is the port on which to serve HTTPS with authentication and authorization
	BindPort int

	APIRequestTimeout int
	logger            logr.Logger

	Done context.Context
}

type APIServer struct {
	engine *gin.Engine
	stopCh context.Context
	logger logr.Logger
	port   int
}

func NewAPIServer(opts NewAPIServerOpts, svc V1Service) APIServer {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logger.SetLogger(logger.WithWriter(os.Stdout)))
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	v1.RegisterHandlersWithOptions(engine, svc, v1.GinServerOptions{
		BaseURL: "/api/v1",
	})

	return APIServer{
		engine: engine,
		port:   opts.BindPort,
		stopCh: opts.Done,
	}
}

func (s *APIServer) Run() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.port),
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(err, "http server stopped")
		}
	}()

	<-s.stopCh.Done()
	s.logger.Info("gracefully shutting down http server!")
	_ = srv.Shutdown(context.Background())
}
