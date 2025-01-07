package resizer

import (
	"log/slog"
	"os"
	"testing"

	"github.com/derticom/image-previewer/internal/model"

	"github.com/stretchr/testify/require"
)

func TestResizer_ResizeImage(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	resizer := New(log)

	data, err := os.ReadFile("./test_images/_gopher_original_1024x504.jpg")
	require.NoError(t, err)

	tests := []struct {
		name               string
		img                *model.Image
		wantResizedImgPath string
		wantErr            bool
	}{
		{
			name: "OK case",
			img: &model.Image{
				Source: "https://raw.githubusercontent.com/_gopher_original_1024x504.jpg",
				Width:  300,
				Height: 200,
				Data:   data,
			},
			wantResizedImgPath: "https:__raw.githubusercontent.com__gopher_original_1024x504.jpg",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResizedImgPath, err := resizer.ResizeImage(tt.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResizeImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResizedImgPath != tt.wantResizedImgPath {
				t.Errorf("ResizeImage() gotResizedImgPath = %v, want %v", gotResizedImgPath, tt.wantResizedImgPath)
			}
		})
	}
}
