package http

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &Server{e: e}
}

func (s *Server) RegisterRoutes(register func(*echo.Group)) {
	api := s.e.Group("/api/v1")
	register(api)
}

func (s *Server) Start(addr string) error {
	return s.e.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) StartGraceful(addr string, shutdownTimeout time.Duration) error {
	go func() {
		if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan struct{}, 1)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := s.e.Shutdown(ctx); err != nil {
			s.e.Logger.Fatalf("server forced to shutdown: %v", err)
		}
	}()

	return nil
}