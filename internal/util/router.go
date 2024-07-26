package util

import (
	"fmt"
	"strings"
)

// Formats Mux route with method
func Rpm(basePath string) func(method string, path ...string) string {
	return func(method string, path ...string) string {
		if len(path) == 0 {
			return fmt.Sprintf("%s /%s", method, basePath)
		}
		return fmt.Sprintf("%s /%s/%s", method, basePath, strings.Join(path, "/"))
	}
}
