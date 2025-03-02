package main

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
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

	flag.IntVar(&cfg.port, "port", getEnvAsInt("PORT", 4000), "Server port to listen on")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	if err := app.serve(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
