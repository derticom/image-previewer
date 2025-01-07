package client

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/derticom/image-previewer/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_DownloadImage(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	client := New(5*time.Second, log)

	type args struct {
		imgURL  string
		headers model.Headers
	}
	tests := []struct {
		name    string
		args    args
		wantImg bool
		wantErr bool
	}{
		{
			name: "OK case",
			args: args{
				imgURL: "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
				headers: model.Headers{
					"Test-Header": []string{"Some value"},
				},
			},
			wantImg: true,
			wantErr: false,
		},
		{
			name: "Bad Request",
			args: args{
				imgURL: "https://raw.githubusercontent.com/gopher_original_1024x504.jpg",
				headers: model.Headers{
					"Test-Header": []string{"Some value"},
				},
			},
			wantImg: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotImg, err := client.DownloadImage(context.Background(), tt.args.imgURL, tt.args.headers)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantImg, gotImg != nil)
		})
	}
}
