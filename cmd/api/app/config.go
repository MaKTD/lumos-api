package app

import (
	"doctormakarhina/lumos/internal/inra/boot"
	"doctormakarhina/lumos/internal/pkg/envconf"
	"doctormakarhina/lumos/internal/pkg/errs"
	"time"
)

type config struct {
	common    boot.CommonConfig
	log       boot.LoggerConfig
	pg        boot.PgConf
	tgBot     tgBotConfig
	http      boot.HttpConfig
	handlers  handlersConf
	unisender Unisender
}

type tgBotConfig struct {
	Token         string        `env:"TG_BOT_TOKEN,required"`
	PollerTimeout time.Duration `env:"TG_BOT_POLLER_TIMEOUT" envDefault:"30s"`
	Debug         bool          `env:"TG_BOT_DEBUG"`
	AdminChatID   int64         `env:"TG_BOT_ADMIN_CHAT_ID,required"`
}

type handlersConf struct {
	PingRoute              string `env:"HTTP_PING_ROUTE" envDefault:"/ping"`
	ApiServePrefix         string `env:"HTTP_API_SERVE_PREFIX" envDefault:"/api"`
	ApiCorsAllowedHosts    string `env:"HTTP_API_CORS_ALLOWED_HOSTS" envDefault:"http://localhost"`
	StaticServePrefix      string `env:"HTTP_STATIC_SERVE_PREFIX" envDefault:"/static/"`
	HtmlServerPrefix       string `env:"HTTP_HTML_SERVE_PREFIX" envDefault:"/"`
	StaticServePath        string `env:"HTTP_STATIC_SOURCE_PATH" envDefault:"./web/assets"`
	TrialPaymentsRouteHash string `env:"HTTP_TRIAL_PAYMENTS_ROUTE_HASH,required"`
}

type Unisender struct {
	ApiKey                     string `env:"UNISENDER_API_KEY,required"`
	AfterTrialExpiredListTitle string `env:"UNISENDER_AFTER_TRIAL_EXPIRED_LIST_TITLE" envDefault:"Lumos закончился пробный"`
}

func (r *config) load() error {
	return errs.First(
		envconf.Load(&r.common),
		envconf.Load(&r.log),
		envconf.Load(&r.pg),
		envconf.Load(&r.tgBot),
		envconf.Load(&r.http),
		envconf.Load(&r.handlers),
		envconf.Load(&r.unisender),
	)
}
