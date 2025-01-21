package utils

import (
	"fmt"
	"strings"
)

func SourceToFileName(width, height int, source string) string {
	var fileName string
	if !strings.HasSuffix(source, ".jpg") && !strings.HasSuffix(source, ".jpeg") {
		fileName = fmt.Sprintf("%d_%d_%s.jpg", width, height, source)
	} else {
		fileName = fmt.Sprintf("%d_%d_%s", width, height, source)
	}

	return strings.ReplaceAll(fileName, "/", "_")
}
