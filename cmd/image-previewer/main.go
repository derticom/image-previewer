package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/derticom/image-previewer/internal/cache"
	"github.com/derticom/image-previewer/internal/client"
	"github.com/derticom/image-previewer/internal/config"
	"github.com/derticom/image-previewer/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	cfg := config.NewConfig()

	logger, err := setupLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to setup loggerr: %+v", err))
	}

	imgCache := cache.New(cfg.CacheSize, logger.WithGroup("cache"))

	downloader := client.New(cfg.Server.Timeout, logger.WithGroup("downloader"))

	proxyServer := server.New(
		imgCache,
		downloader,
		cfg.Server.Address,
		cfg.Server.Timeout,
		logger.WithGroup("server"),
	)

	err = proxyServer.Run(ctx)
	if err != nil {
		logger.Error("failed to proxyServer.Run", "error", err)
	}
}

func setupLogger(cfg *config.Config) (*slog.Logger, error) {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		),
	)

	return logger, nil
}
