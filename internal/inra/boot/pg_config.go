package boot

import "time"

type PgConf struct {
	Url               string        `env:"PG_URL,required"`
	MaxConns          int         `env:"PG_MAX_CONN" envDefault:"20"`
	MaxConnIdleTime   time.Duration `env:"PG_MAX_CONN_IDLE_TIME" envDefault:"30s"`
}
