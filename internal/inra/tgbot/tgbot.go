package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

const telegramMaxMessageLen = 4096

type Bot struct {
	token         string
	chatID        int64
	debug         bool
	pollerTimeout time.Duration

	logger *slog.Logger
	bot    *tele.Bot
}

type BotCfg struct {
	Token         string
	ChatID        int64
	Debug         bool
	PollerTimeout time.Duration
	Logger        *slog.Logger
}

func NewAdminTgBot(cfg BotCfg) (*Bot, error) {
	logger := cfg.Logger.With("context", "AdminTgBot")

	bot := &Bot{
		token:         cfg.Token,
		chatID:        cfg.ChatID,
		debug:         cfg.Debug,
		pollerTimeout: cfg.PollerTimeout,

		logger: logger,
	}

	err := bot.Init()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (r *Bot) Init() error {
	settings := tele.Settings{
		Token:   r.token,
		Poller:  &tele.LongPoller{Timeout: r.pollerTimeout},
		Verbose: r.debug,
		OnError: func(err error, _ tele.Context) {
			r.logger.Error(
				"unexpected error on admin tg bot",
				slog.String("err", err.Error()),
			)
		},
	}
	bot, err := tele.NewBot(settings)
	if err != nil {
		return err
	}
	r.bot = bot

	r.bot.Use(ChatWhitelist(r.chatID))

	r.bot.Handle("/start", r.HandleOnStart)

	return nil
}

func (r *Bot) Run(ctx context.Context) error {
	if r.bot == nil {
		panic("error: admin tg bot is not initialized")
	}

	stopped := make(chan struct{})

	go func() {
		defer close(stopped)
		r.bot.Start()
	}()

	for {
		select {
		case <-stopped:
			return fmt.Errorf("tgbot stopped unexpectedly")
		case <-ctx.Done():
			r.bot.Stop()
			<-stopped
			return nil
		}
	}
}

func (r *Bot) HandleOnStart(ctx tele.Context) error {
	return ctx.Send(welcomeText)
}
