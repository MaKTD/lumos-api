package boot

import (
	"time"
)

type HttpConfig struct {
	Addr              string        `env:"HTTP_ADDR" envDefault:":8080"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	MaxHeaderBytes    int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"`
	ShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"30s"`
	AutoCertEnabled   bool          `env:"HTTP_AUTO_CERT_ENABLED" envDefault:"false"`
	AutoCertCacheDir  string        `env:"HTTP_AUTO_CERT_CACHE_DIR" envDefault:"/certificates"`
	AutoCertEmail     string        `env:"HTTP_AUTO_CERT_EMAIL" envDefault:"local@mail.com"`
	AutoCertHosts     string        `env:"HTTP_AUTO_CERT_HOSTS" envDefault:"localhost"`
}
