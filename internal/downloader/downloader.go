// Package downloader предназначен для скачивания исходного изображения из сети.
package downloader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/derticom/image-previewer/internal/model"

	"github.com/pkg/errors"
)

const methodType = "GET"

type Downloader struct {
	client *http.Client
	log    *slog.Logger
}

func New(timeout time.Duration, log *slog.Logger) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

// DownloadImage принимает входящий запрос, обрабатывает его и скачивает изображение.
// Возвращает объект model.Image, содержащий изображение в байтовом формате и его название.
func (c *Downloader) DownloadImage(
	ctx context.Context,
	request model.Request,
) (img *model.Image, err error) {
	if !strings.HasPrefix(request.URL, "http") {
		request.URL = "https://" + request.URL
	}

	url := request.URL
	if request.Params != "" {
		url = fmt.Sprintf("%s?%s", request.URL, request.Params)
	}

	req, err := http.NewRequestWithContext(ctx, methodType, url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to http.NewRequest")
	}

	for key, values := range request.Headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to client.Do")
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			c.log.Error("failed to Body.Close", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		c.log.Error("unexpected status code while downloading image",
			"URL", request.URL,
			"StatusCode", resp.StatusCode,
			"Status", resp.Status)
		return nil, errors.New("unexpected status code")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to io.ReadAll")
	}

	return &model.Image{
		Source: request.URL,
		Data:   data,
	}, nil
}
