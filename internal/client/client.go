package client

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/derticom/image-previewer/internal/model"

	"github.com/pkg/errors"
)

const methodType = "GET"

type Client struct {
	client *http.Client
	log    *slog.Logger
}

func New(timeout time.Duration, log *slog.Logger) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

func (c *Client) DownloadImage(
	ctx context.Context,
	imgURL string,
	requestHeaders model.Headers,
) (img *model.Image, err error) {
	if !strings.HasPrefix(imgURL, "http") {
		imgURL = "https://" + imgURL
	}

	req, err := http.NewRequestWithContext(ctx, methodType, imgURL, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to http.NewRequest")
	}

	for key, values := range requestHeaders {
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
			"URL", imgURL,
			"StatusCode", resp.StatusCode,
			"Status", resp.Status)
		return nil, errors.New("unexpected status code")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to io.ReadAll")
	}

	return &model.Image{
		Source: imgURL,
		Data:   data,
	}, nil
}
