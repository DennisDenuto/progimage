package apiserver

import (
	"github.com/go-logr/logr"
	"net/http"
)

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
}

func (s *APIServer) Run() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.logger.Error(err, "API server stopped")
		}
	}()
}

func NewAPIServer(opts NewAPIServerOpts) APIServer {
	return APIServer{}
}
