package utils

import (
	"fmt"
	"strings"
)

func SourceToFileName(width, height int, source string) string {
	name := fmt.Sprintf("%d_%d_%s", width, height, source)
	return strings.ReplaceAll(name, "/", "_")
}
