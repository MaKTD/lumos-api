package boot

import (
	"context"
	"database/sql"
	"doctormakarhina/lumos/internal/api"
	"doctormakarhina/lumos/internal/html"
	"doctormakarhina/lumos/internal/pkg/configs"
	"doctormakarhina/lumos/internal/pkg/errs"
	"doctormakarhina/lumos/internal/pkg/httpx"
	"doctormakarhina/lumos/internal/pkg/logger"
	"doctormakarhina/lumos/internal/pkg/modules/auth_logs"
	"doctormakarhina/lumos/internal/pkg/modules/search_logs"
	"doctormakarhina/lumos/internal/pkg/postgresql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
)

type App struct {
	OsCtx       context.Context
	OsCtxCancel context.CancelFunc

	RootLogger *slog.Logger

	Conf *ConfRegistry

	DbPool *sql.DB

	HttpRouter *chi.Mux
	HttpServer *httpx.Server

	authLogsService   *auth_logs.Service
	searchLogsService *search_logs.Service
}

func (a *App) Init() error {
	a.OsCtx, a.OsCtxCancel = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	logger.ConfigureDefault()

	a.Conf = &ConfRegistry{}
	configs.LoadDotenvIfEnabled()
	err := a.Conf.Load()
	if err != nil {
		return err
	}

	a.RootLogger, err = logger.NewRoot(logger.Config(a.Conf.Log))
	if err != nil {
		return err
	}

	a.DbPool, err = postgresql.NewPgPool(a.OsCtx, postgresql.Config(a.Conf.Pg), a.RootLogger)
	if err != nil {
		return err
	}

	a.authLogsService = auth_logs.NewService(a.DbPool)
	a.searchLogsService = search_logs.NewService(a.DbPool)

	a.HttpRouter = chi.NewRouter()
	a.HttpRouter.Use(
		middleware.Recoverer,
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: strings.Split(a.Conf.AppHttpHandle.ApiCorsAllowedHosts, ","),
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	a.HttpRouter.Route(a.Conf.AppHttpHandle.ApiServePrefix, func(r chi.Router) {
		r.Get(a.Conf.AppHttpHandle.PingRoute, api.NewPingHandlerFunc())
		api.NewGlobalHandler(a.RootLogger, a.authLogsService, a.searchLogsService).RegIn(r)
	})
	a.HttpRouter.Route(a.Conf.AppHttpHandle.HtmlServerPrefix, func(r chi.Router) {
		html.NewGlobalHandler(a.RootLogger).RegIn(r)
	})

	// TODO(add client side caching, etag probably??)
	a.HttpRouter.Handle(
		fmt.Sprintf("%s*", a.Conf.AppHttpHandle.StaticServePrefix),
		http.StripPrefix(
			a.Conf.AppHttpHandle.StaticServePrefix,
			http.FileServer(http.Dir(a.Conf.AppHttpHandle.StaticServePath)),
		),
	)

	a.HttpServer = httpx.NewServer(
		httpx.Cfg(a.Conf.AppHttpServer),
		a.HttpRouter,
		a.RootLogger,
	)

	return nil
}

func (a *App) Run() error {
	err := a.DbPool.PingContext(a.OsCtx)
	if err != nil {
		return errs.WrapErrorf(err, errs.ErrCodeInternal, "initial db pool ping failed")
	}

	doneHttpServerCtx := a.HttpServer.GoStart()

	select {
	case <-a.OsCtx.Done():
		a.RootLogger.Info("received close signal from os, shutting down")
		return nil
	case <-doneHttpServerCtx.Done():
		return errs.WrapErrorf(doneHttpServerCtx.Err(), errs.ErrCodeInternal, "http server closed unexpectedly, shutting down")
	}
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), a.Conf.Boot.ShutdownTimeout)
	defer cancel()

	if a.HttpServer != nil {
		err := a.HttpServer.Shutdown(ctx)
		if err != nil {
			a.RootLogger.Error("http server failed to shutdown properly", slog.String("err", err.Error()))
		} else {
			a.RootLogger.Info("http server successfully shutdown")
		}
	}

	if a.DbPool != nil {
		err := a.DbPool.Close()
		if err != nil {
			a.RootLogger.Error("db pool failed to shutdown properly", slog.String("err", err.Error()))
		} else {
			a.RootLogger.Info("db pool successfully closed")
		}
	}
}

func StartApp() error {
	a := App{}
	defer a.Shutdown()

	err := a.Init()
	if err != nil {
		return err
	}

	return a.Run()
}
