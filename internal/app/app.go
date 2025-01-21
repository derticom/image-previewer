package app

import (
	"context"
	"log/slog"

	"github.com/derticom/image-previewer/config"
	"github.com/derticom/image-previewer/internal/cache"
	"github.com/derticom/image-previewer/internal/downloader"
	"github.com/derticom/image-previewer/internal/server"

	"github.com/pkg/errors"
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	imgCache := cache.New(cfg.CacheSize, log.WithGroup("cache"))

	imgDownloader := downloader.New(cfg.Server.Timeout, log.WithGroup("downloader"))

	proxyServer := server.New(
		imgCache,
		imgDownloader,
		cfg.Server.Address,
		cfg.Server.Timeout,
		log.WithGroup("server"),
	)

	err := proxyServer.Run(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to proxyServer.Run")
	}

	return nil
}
