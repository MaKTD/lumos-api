package app

import (
	"doctormakarhina/lumos/internal/inra/boot"
	"doctormakarhina/lumos/internal/pkg/envconf"
	"doctormakarhina/lumos/internal/pkg/errs"
	"time"
)

type config struct {
	common   boot.CommonConfig
	log      boot.LoggerConfig
	pg       boot.PgConf
	tgBot    tgBotConfig
	http     boot.HttpConfig
	handlers handlersConf
}

type tgBotConfig struct {
	Token         string        `env:"TG_BOT_TOKEN,required"`
	PollerTimeout time.Duration `env:"TG_BOT_POLLER_TIMEOUT" envDefault:"30s"`
	Debug         bool          `env:"TG_BOT_DEBUG"`
	// AdminChatID   int64         `env:"ADMIN_TG_BOT_ADMIN_CHAT_ID,required"`
}

type handlersConf struct {
	PingRoute           string `env:"HTTP_PING_ROUTE" envDefault:"/ping"`
	ApiServePrefix      string `env:"HTTP_API_SERVE_PREFIX" envDefault:"/api"`
	ApiCorsAllowedHosts string `env:"HTTP_API_CORS_ALLOWED_HOSTS" envDefault:"http://localhost"`
	StaticServePrefix   string `env:"HTTP_STATIC_SERVE_PREFIX" envDefault:"/static/"`
	HtmlServerPrefix    string `env:"HTTP_HTML_SERVE_PREFIX" envDefault:"/"`
	StaticServePath     string `env:"HTTP_STATIC_SOURCE_PATH" envDefault:"./web/assets"`
}

func (r *config) load() error {
	return errs.First(
		envconf.Load(&r.common),
		envconf.Load(&r.log),
		envconf.Load(&r.pg),
		envconf.Load(&r.tgBot),
		envconf.Load(&r.http),
		envconf.Load(&r.handlers),
	)
}
