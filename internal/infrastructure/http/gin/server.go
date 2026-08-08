package gin

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	router := NewRouter()

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
