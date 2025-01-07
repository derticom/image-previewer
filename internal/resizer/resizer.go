package resizer

import (
	"image"
	"image/color"
	"log/slog"
	"os"

	"github.com/derticom/image-previewer/internal/model"
	"github.com/derticom/image-previewer/internal/utils"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

const (
	tmpPath     = "./tmp"
	StoragePath = "./storage/"
)

type Resizer struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Resizer {
	return &Resizer{log: log}
}

func (r *Resizer) ResizeImage(img *model.Image) (resizedImgPath string, err error) {
	tempFile, err := os.CreateTemp(tmpPath, "tmp-*.jpg")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		errRemove := os.Remove(name)
		if errRemove != nil {
			r.log.Error("failed to remove temporary file", "error", errRemove)
		}
	}(tempFile.Name())

	_, err = tempFile.Write(img.Data)
	if err != nil {
		return "", errors.Wrap(err, "failed to tempFile.Write")
	}
	defer func(tempFile *os.File) {
		errClose := tempFile.Close()
		if errClose != nil {
			r.log.Error("failed to close tempFile", "error", errClose)
		}
	}(tempFile)

	fileName := utils.SourceToFileName(img.Width, img.Height, img.Source)
	src, err := imaging.Open(tempFile.Name())
	if err != nil {
		return "", errors.Wrap(err, "failed to imaging.Open")
	}

	resizedImg := imaging.Fill(src, img.Width, img.Height, imaging.Center, imaging.Lanczos)

	// Create a new image and paste the four produced images into it.
	dst := imaging.New(img.Width, img.Height, color.NRGBA{})
	dst = imaging.Paste(dst, resizedImg, image.Pt(0, 0))

	// Save the resulting image as JPEG.

	err = imaging.Save(dst, StoragePath+fileName)
	if err != nil {
		return "", errors.Wrap(err, "failed to imaging.Save")
	}

	return fileName, nil
}
