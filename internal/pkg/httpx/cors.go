package httpx

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type CorsConfig struct {
	AllowedOrigins       string
	AllowedMethods       string
	AllowedHeaders       string
	ExposedHeaders       string
	MaxAge               int
	AllowCredentials     bool
	AllowPrivateNetwork  bool
	OptionsPassthrough   bool
	OptionsSuccessStatus int
	Debug                bool
}

func NewCordsHandler(config CorsConfig, handler http.Handler) http.Handler {
	if config.AllowedOrigins == "" {
		config.AllowedOrigins = "*"
	}
	if config.AllowedMethods == "" {
		config.AllowedOrigins = "HEAD,GET,POST"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:       strings.Split(strings.TrimSpace(config.AllowedOrigins), ","),
		AllowedMethods:       strings.Split(strings.TrimSpace(config.AllowedMethods), ","),
		AllowedHeaders:       strings.Split(strings.TrimSpace(config.AllowedHeaders), ","),
		ExposedHeaders:       strings.Split(strings.TrimSpace(config.ExposedHeaders), ","),
		MaxAge:               config.MaxAge,
		AllowCredentials:     config.AllowCredentials,
		AllowPrivateNetwork:  config.AllowPrivateNetwork,
		OptionsPassthrough:   config.OptionsPassthrough,
		OptionsSuccessStatus: config.OptionsSuccessStatus,
		Debug:                config.Debug,
	})

	return c.Handler(handler)
}
