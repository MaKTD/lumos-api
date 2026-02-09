package boot

import "time"

type CommonConfig struct {
	AppName         string        `env:"APP_NAME" envDefault:"APP"`
	ShutdownTimeout time.Duration `env:"BOOTSTRAP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}
