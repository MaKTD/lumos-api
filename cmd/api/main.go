package main

import (
	"doctormakarhina/lumos/internal/boot"
	"log/slog"
	"os"
)

func main() {
	err := boot.StartApp()
	if err != nil {
		slog.Error("application bootstrap failed with error", slog.String("err", err.Error()))
		os.Exit(1)
	} else {
		slog.Info("application was successfully closed")
		os.Exit(0)
	}
}
