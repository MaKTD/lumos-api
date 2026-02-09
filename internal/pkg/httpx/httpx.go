package httpx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	shutdownTimeout time.Duration
	logger          *slog.Logger
	server          *http.Server
}

type Config struct {
	Name              string
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	ShutdownTimeout   time.Duration
	TLSconfig         *tls.Config
	AutoCertEnabled   bool
	AutoCertCacheDir  string
	AutoCertEmail     string
	AutoCertHosts     string
}

func NewServer(config Config, handler http.Handler, rootLogger *slog.Logger) *Server {
	if config.Name == "" {
		config.Name = "Server"
	}
	if config.Addr == "" {
		config.Addr = ":8080"
	}
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 5 * time.Second
	}

	var logger *slog.Logger
	if rootLogger != nil {
		logger = rootLogger.With(slog.String("context", fmt.Sprintf("Http%s", config.Name)))
	} else {
		logger = slog.Default()
	}

	server := &Server{
		shutdownTimeout: config.ShutdownTimeout,
		logger:          logger,
		server: &http.Server{
			Addr:              config.Addr,
			Handler:           handler,
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
			ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
			TLSConfig:         config.TLSconfig,
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

func (r *Server) Run(ctx context.Context) error {
	srvCtx, srvCancel := context.WithCancelCause(context.Background())

	go func() {
		srvCancel(r.start())
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), r.shutdownTimeout)
		defer shutdownCancel()
		err := r.shutdown(shutdownCtx)
		<-srvCtx.Done()
		return err
	case <-srvCtx.Done():
		return srvCtx.Err()
	}
}

func (r *Server) start() error {
	r.logger.Info(fmt.Sprintf("server is starting to listen on addr = %s", r.server.Addr))

	if r.server.TLSConfig != nil {
		if err := r.server.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
			r.logger.Error("server closed unexpectedly", slog.String("err", err.Error()))
			return err
		}
	} else {
		if err := r.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			r.logger.Error("server closed unexpectedly", slog.String("err", err.Error()))
			return err
		}
	}

	r.logger.Info("server was successfully closed")
	return nil
}

func (r *Server) shutdown(ctx context.Context) error {
	r.logger.Info("start shutting down server")
	err := r.server.Shutdown(ctx)
	if err != nil {
		r.logger.Error("shutdown server error", slog.String("err", err.Error()))
	}

	return err
}
