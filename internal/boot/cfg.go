package boot

import (
	"doctormakarhina/lumos/internal/pkg/configs"
	"doctormakarhina/lumos/internal/pkg/errs"
	"time"
)

type ConfRegistry struct {
	App           AppConf
	Boot          BootConf
	Pg            PgConf
	Log           LogConf
	AppHttpServer AppHttpServerConf
	AppHttpHandle AppHttpHandleConf
}

func (r *ConfRegistry) Load() error {
	return errs.First(
		configs.Load(&r.App),
		configs.Load(&r.Boot),
		configs.Load(&r.Log),
		configs.Load(&r.Pg),
		configs.Load(&r.AppHttpServer),
		configs.Load(&r.AppHttpHandle),
	)
}

type BootConf struct {
	BootTimeout     time.Duration `env:"BOOSTRAP_BOOT_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout time.Duration `env:"BOOSTRAP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type AppConf struct {
	Name string `env:"APP_NAME,required"`
}

type PgConf struct {
	Url               string        `env:"PG_URL,required"`
	MaxConns          int32         `env:"PG_MAX_CONN" envDefault:"20"`
	MinConns          int32         `env:"PG_MIN_CONN" envDefault:"1"`
	HealthCheckPeriod time.Duration `env:"PG_HEALTH_CHECK_PERIOD" envDefault:"5s"`
	MaxConnIdleTime   time.Duration `env:"PG_MAX_CONN_IDLE_TIME" envDefault:"30s"`
	Debug             bool          `env:"PG_DEBUG" envDefault:"false"`
	ConnectionTimeout time.Duration `env:"PG_CONNECTION_TIMEOUT" envDefault:"5s"`
}

type LogConf struct {
	Level          string `env:"LOG_LEVEL" envDefault:"info"`
	Pretty         bool   `env:"LOG_PRETTY"`
	IncludeSources bool   `env:"LOG_SOURCES"`
}

type AppHttpServerConf struct {
	Name              string        `env:"HTTP_SERVER_NAME" envDefault:"httpApiServer"`
	ListenPort        int           `env:"HTTP_LISTEN_PORT" envDefault:"8080"`
	ListenHost        string        `env:"HTTP_LISTEN_HOST" envDefault:"127.0.0.1"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"0"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"0"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"5m"`
}

type AppHttpHandleConf struct {
	PingRoute         string `env:"HTTP_PING_ROUTE" envDefault:"/ping"`
	ApiServePrefix    string `env:"HTTP_API_SERVE_PREFIX" envDefault:"/api"`
	StaticServePrefix string `env:"HTTP_STATIC_SERVE_PREFIX" envDefault:"/static/"`
	HtmlServerPrefix  string `env:"HTTP_HTML_SERVE_PREFIX" envDefault:"/"`
	StaticServePath   string `env:"HTTP_STATIC_SOURCE_PATH" envDefault:"./web/assets"`
}
