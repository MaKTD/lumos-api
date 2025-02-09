package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	server http.Server
	logger *slog.Logger
}

type Cfg struct {
	Name              string
	ListenPort        int
	ListenHost        string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("server is starting to listen on addr = %s", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("server closed unexpectedly", slog.String("err", err.Error()))
		return err
	}
	s.logger.Info("server was successfully closed")
	return nil
}

func (s *Server) GoStart() context.Context {
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		cancel(s.Start())
	}()

	return ctx
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("start shutting down")
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("shutdown server error", slog.String("err", err.Error()))
	}
	return err
}

func NewServer(config Cfg, handler http.Handler, pLogger *slog.Logger) *Server {
	logger := pLogger.With(slog.String("context", config.Name))

	server := &Server{
		logger: logger,
		server: http.Server{
			Addr:                         fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort),
			Handler:                      handler,
			DisableGeneralOptionsHandler: false,
			ReadHeaderTimeout:            config.ReadHeaderTimeout,
			WriteTimeout:                 config.WriteTimeout,
			IdleTimeout:                  config.IdleTimeout,
		},
	}

	return server
}
