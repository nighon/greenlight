package main

import (
	"log/slog"
	"os"
	"sync"
)

type config struct {
	port int
}

type application struct {
	config config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg.port = 4000

	app := &application{
		config: cfg,
		logger: logger,
	}

	if err := app.serve(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
