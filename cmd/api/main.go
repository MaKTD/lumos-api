package main

import (
	"doctormakarhina/lumos/cmd/api/app"
	"doctormakarhina/lumos/internal/inra/boot"
	"log/slog"
	"os"
)

func main() {
	err := boot.StartApp(&app.Server{})
	if err != nil {
		slog.Error("error during bootstrap or running app", slog.String("err", err.Error()))
		os.Exit(1)
	} else {
		slog.Info("app has been shutdowned")
	}
}
