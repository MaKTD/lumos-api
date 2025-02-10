package httpx

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	config Cfg
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
	AutoCertEnabled   bool
	AutoCertCacheDir  string
	AutoCertEmail     string
	AutoCertHosts     string
}

func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("server is starting to listen on addr = %s", s.server.Addr))

	if s.config.AutoCertEnabled {
		if err := s.server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server closed unexpectedly", slog.String("err", err.Error()))
			return err
		}
		s.logger.Info("server was successfully closed")
		return nil
	} else {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server closed unexpectedly", slog.String("err", err.Error()))
			return err
		}
		s.logger.Info("server was successfully closed")
		return nil
	}
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
		config: config,
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

	if config.AutoCertEnabled {
		hosts := strings.Split(config.AutoCertHosts, ",")

		manager := &autocert.Manager{
			Cache:      autocert.DirCache(config.AutoCertCacheDir),
			Prompt:     autocert.AcceptTOS,
			Email:      config.AutoCertEmail,
			HostPolicy: autocert.HostWhitelist(hosts...),
		}

		server.server.TLSConfig = manager.TLSConfig()
	}

	return server
}
