package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"example.com/internal/data"
	"example.com/internal/vcs"
	_ "github.com/go-sql-driver/mysql"
)

var (
	version = vcs.Version()
)

type config struct {
	port int
	db   struct {
		dsn string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	version = getEnvAsString("VERSION", version)

	flag.IntVar(&cfg.port, "port", getEnvAsInt("PORT", 4000), "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "db-dsn", getEnvAsString("DATABASE_URL", ""), "MySQL DSN")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", getEnvAsFloat64("LIMITER_RPS", 2), "Rate limit to apply to requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", getEnvAsInt("LIMITER_BURST", 4), "Burst limit to apply to requests")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", getEnvAsBool("LIMITER_ENABLED", true), "Enable rate limiting")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: \t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error("error opening db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	if err := app.serve(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create context with 5s timeout. If we can't connect in 5s, we cancel the context and return an error.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func getEnvAsString(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvAsFloat64(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
