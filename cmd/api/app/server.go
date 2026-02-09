package app

import (
	"context"
	"doctormakarhina/lumos/internal/core/payments"
	"doctormakarhina/lumos/internal/inra/emails"
	"doctormakarhina/lumos/internal/inra/httpapi"
	"doctormakarhina/lumos/internal/inra/pg"
	"doctormakarhina/lumos/internal/inra/tgbot"
	"doctormakarhina/lumos/internal/pkg/db"
	"doctormakarhina/lumos/internal/pkg/envconf"
	"doctormakarhina/lumos/internal/pkg/httpx"
	"doctormakarhina/lumos/internal/pkg/logger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/run"
)

type Server struct {
	cfg        *config
	rootLogger *slog.Logger
	db         *sqlx.DB
	bot        *tgbot.Bot
	api        *httpx.Server
}

func (r *Server) Init() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	logger.ConfigureDefault()

	envconf.LoadDotenvIfEnabled()
	r.cfg = &config{}
	err := r.cfg.load()
	if err != nil {
		return err
	}

	r.rootLogger, err = logger.New(
		r.cfg.log.SlogLevel(),
		r.cfg.log.Pretty,
		r.cfg.log.IncludeSources,
		slog.String("APP_NAME", r.cfg.common.AppName),
	)
	if err != nil {
		return err
	}
	logger.SetToDefault(r.rootLogger)

	r.db, err = db.NewPG(
		ctx,
		r.cfg.pg.Url,
		r.cfg.pg.MaxConns,
		r.cfg.pg.MaxConns,
		r.cfg.pg.MaxConnIdleTime,
	)
	if err != nil {
		return err
	}
	err = r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	r.bot, err = tgbot.NewAdminTgBot(tgbot.BotCfg{
		Token:         r.cfg.tgBot.Token,
		ChatID:        r.cfg.tgBot.AdminChatID,
		Debug:         r.cfg.tgBot.Debug,
		PollerTimeout: r.cfg.tgBot.PollerTimeout,
		Logger:        r.rootLogger,
	})

	usersRepo := pg.NewUserRepo(r.db)

	emailSrv := emails.NewUniSenderSrv(r.cfg.unisender.ApiKey)

	paymentSrv := payments.NewPaymentsService(
		usersRepo,
		emailSrv,
		r.bot,
	)

	apiHandler := httpapi.NewRouter()
	apiHandler.Use(
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedOrigins:   strings.Split(r.cfg.handlers.ApiCorsAllowedHosts, ","),
			AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH", "HEAD"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			// Maximum value not ignored by any of major browsers
			MaxAge: 3600,
		}),
	)
	apiHandler.Route(r.cfg.handlers.ApiServePrefix, func(router chi.Router) {
		httpapi.RegInPing(router)
		httpapi.RegInHealthz(router, r.rootLogger)
		httpapi.RegInAuthLogs(router, r.db, r.rootLogger)
		httpapi.RegInSearchLogs(router, r.db, r.rootLogger)
		httpapi.RegInTrialPayments(
			router,
			r.cfg.handlers.TrialPaymentsRouteHash,
			paymentSrv,
			r.bot,
		)
	})

	// TODO(add client side caching, etag probably??)
	apiHandler.Handle(
		fmt.Sprintf("%s*", r.cfg.handlers.StaticServePrefix),
		http.StripPrefix(
			r.cfg.handlers.StaticServePrefix,
			http.FileServer(http.Dir(r.cfg.handlers.StaticServePath)),
		),
	)

	r.api = httpx.NewServer(httpx.Config{
		Name:              "Api",
		Addr:              r.cfg.http.Addr,
		ReadTimeout:       r.cfg.http.ReadTimeout,
		ReadHeaderTimeout: r.cfg.http.ReadHeaderTimeout,
		WriteTimeout:      r.cfg.http.WriteTimeout,
		IdleTimeout:       r.cfg.http.IdleTimeout,
		MaxHeaderBytes:    r.cfg.http.MaxHeaderBytes,
		ShutdownTimeout:   r.cfg.http.ShutdownTimeout,
	}, apiHandler, r.rootLogger)

	r.rootLogger.Info("app initialized")

	return nil
}

func (r *Server) Run() error {
	var g run.Group

	{
		logger := r.rootLogger.With(slog.String("context", "SignalListener"))
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case sig := <-term:
					logger.Info("Received as OS signal, exiting gracefully...", slog.String("signal", sig.String()))
				case <-cancel:
				}
				return nil
			},
			func(err error) {
				close(cancel)
			},
		)
	}
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(
			func() error {
				return r.bot.Run(ctx)
			},
			func(err error) {
				cancel()
			},
		)
	}
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(
			func() error {
				return r.api.Run(ctx)
			},
			func(err error) {
				cancel()
			},
		)
	}

	return g.Run()
}

func (r *Server) Shutdown() {
	if r.db != nil {
		err := r.db.Close()
		if err != nil && r.rootLogger != nil {
			r.rootLogger.Error(
				"failed to gracefully close db",
				slog.String("err", err.Error()),
			)
		}
	}
}
