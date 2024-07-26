package util

import (
	"errors"
	"os"
	"strings"
)

func FileExists(path string) (bool, error) {
	v, err := os.Stat(path)
	if err == nil && !v.IsDir() {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func FileWritable(path string) (bool, error) {
	v, err := os.OpenFile(path, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return false, err
	}
	v.Close()

	return true, nil
}

func FileHasExt(pathOrName, ext string, otherExts ...string) bool {
	if strings.HasSuffix(pathOrName, ext) {
		return true
	}
	for _, ext := range otherExts {
		if strings.HasSuffix(pathOrName, ext) {
			return true
		}
	}
	return false
}
