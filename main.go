package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/jessestricker/terraform-backend-age/internal"
)

func main() {
	if debug, err := strconv.ParseBool(os.Getenv("TF_BACKEND_AGE_DEBUG")); err == nil && debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	os.Exit(internal.Main())
}
