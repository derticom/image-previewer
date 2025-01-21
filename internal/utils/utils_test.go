package utils

import "testing"

func TestSourceToFileName(t *testing.T) {
	type args struct {
		width  int
		height int
		source string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK case",
			args: args{
				width:  300,
				height: 200,
				source: "https://raw.githubusercontent.com/_gopher_original_1024x504.jpg",
			},
			want: "300_200_https:__raw.githubusercontent.com__gopher_original_1024x504.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SourceToFileName(tt.args.width, tt.args.height, tt.args.source); got != tt.want {
				t.Errorf("SourceToFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
