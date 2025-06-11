package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
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

// AskPass 询问密码
func AskPass(msg string) string {
	fmt.Print(msg)
	pass, err := term.ReadPassword(0)
	if err != nil {
		log.Fatal("读取密码失败:", err)
	}
	fmt.Println("")
	return strings.TrimSpace(string(pass))
}
