package server

import (
	"fmt"
	"net/http"
	"time"

	"FreeMusic/internal/config"

	"golang.org/x/net/context"
)

// server ...
type server struct {
	httpServer *http.Server
	config     config.Config
}

// NewServer ...
func NewServer(conf *config.Config) *server {
	return &server{
		httpServer: nil,
		config:     *conf,
	}
}

// Run ...
func (s *server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.AppHost, s.config.AppPort),
		Handler:        handler,
		MaxHeaderBytes: s.config.AppMaxHeaderBytes, // 1 MB
		ReadTimeout:    time.Duration(s.config.AppReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.config.AppWriteTimeout) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown ...
func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
