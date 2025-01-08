package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/derticom/image-previewer/internal/cache"
	"github.com/derticom/image-previewer/internal/model"
	"github.com/derticom/image-previewer/internal/resizer"
	"github.com/derticom/image-previewer/internal/utils"

	"github.com/pkg/errors"
)

type ImageDownloader interface {
	DownloadImage(ctx context.Context, imgURL string, requestHeaders model.Headers) (img *model.Image, err error)
}

type ImageResizer interface {
	ResizeImage(img *model.Image) (resizedImgPath string, err error)
}

type Cache interface {
	Set(key cache.Key, value any) bool
	Get(key cache.Key) (any, bool)
}

type Server struct {
	cache      Cache
	downloader ImageDownloader
	resizer    ImageResizer

	address string
	timeout time.Duration

	log *slog.Logger
}

func New(
	imgCache Cache,
	downloader ImageDownloader,
	address string,
	timeout time.Duration,
	log *slog.Logger,
) *Server {
	return &Server{
		cache:      imgCache,
		downloader: downloader,
		resizer:    resizer.New(log),
		address:    address,
		timeout:    timeout,
		log:        log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /fill/{width}/{height}/{url...}", s.handle)

	server := http.Server{
		Addr:        s.address,
		Handler:     mux,
		ReadTimeout: s.timeout,
	}

	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil {
			s.log.Error("failed to server.Shutdown", "error", err)
		}
	}()

	s.log.Info("server listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("shutdown server")
			return nil
		}
		return errors.Wrap(err, "failed to server.ListenAndServe")
	}

	return nil
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	s.log.Info("start processing", "request", r.URL.Path)

	width, err := strconv.Atoi(r.PathValue("width"))
	if err != nil {
		s.log.Error("failed process width", "error", err)
		_, err := fmt.Fprintf(w, "failed process width: %s", r.PathValue("width"))
		if err != nil {
			s.log.Error("failed to write response", "error", err)
		}
		return
	}
	height, err := strconv.Atoi(r.PathValue("height"))
	if err != nil {
		s.log.Error("failed process height", "error", err)
		_, err := fmt.Fprintf(w, "failed process height: %s", r.PathValue("height"))
		if err != nil {
			s.log.Error("failed to write response", "error", err)
		}
		return
	}

	url := r.PathValue("url")

	requestHeaders := model.Headers{}
	for key, values := range r.Header {
		requestHeaders[key] = append(requestHeaders[key], values...)
	}

	resizedImageName, ok := s.cache.Get(
		cache.Key(utils.SourceToFileName(width, height, url)),
	)
	if !ok {
		s.log.Info("image not found in cache",
			"url", url,
			"width", width,
			"height", height,
		)

		img, err := s.downloader.DownloadImage(r.Context(), url, requestHeaders)
		if err != nil {
			s.log.Error("failed to download image", "error", err)
			_, err := fmt.Fprintf(w, "failed to download image")
			if err != nil {
				s.log.Error("failed to write response", "error", err)
			}
			return
		}

		img.Width = width
		img.Height = height

		resizedImageName, err = s.resizer.ResizeImage(img)
		if err != nil {
			s.log.Error("failed to resize image", "error", err)
			_, err := fmt.Fprintf(w, "failed to resize image")
			if err != nil {
				s.log.Error("failed to write response", "error", err)
			}
			return
		}

		s.cache.Set(
			cache.Key(utils.SourceToFileName(width, height, url)),
			resizedImageName,
		)
	} else {
		s.log.Info("image found in cache")
	}

	resizedImage, err := os.ReadFile(resizer.StoragePath + resizedImageName.(string))
	if err != nil {
		s.log.Error("failed to read resized image", "error", err)
		_, err := fmt.Fprintf(w, "failed to read resized image")
		if err != nil {
			s.log.Error("failed to write response", "error", err)
		}
		return
	}

	_, err = w.Write(resizedImage)
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		_, err := fmt.Fprintf(w, "failed to write response")
		if err != nil {
			s.log.Error("failed to write response", "error", err)
		}
		return
	}

	s.log.Info("successfully finished processing", "request", r.URL.Path)
}
