package main

import (
	"os"
	"path/filepath"
	"strings"
)

// NormalizePath 规范化路径
func NormalizePath(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}

		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, path[2:])
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}

	cleanPath := filepath.Clean(absPath)

	_, err = os.Stat(cleanPath)
	if err != nil && !os.IsNotExist(err) {
		return cleanPath, err
	}

	return cleanPath, nil
}
