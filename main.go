package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/melbahja/goph"
)

var (
	configPath string = "config.toml"
	config     *Config
	gophConfig *goph.Config
	gophClient *goph.Client
)

func main() {
	var err error
	// 如果config.toml不存在，直接执行git命令
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		pwd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}

		cmd := exec.Command("git", os.Args[1:]...)
		cmd.Dir = pwd
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			// 命令执行失败时返回错误码
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			os.Exit(1)
		}
		return
	}

	config, err = LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config file %s: %v", configPath, err)
	}

	gophConfig, err = GetHostConfig(config.SSH.Host)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	gophClient, err = goph.NewConn(gophConfig)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer gophClient.Close()

	execCmd := "cd " + config.Repo.Path + " && " + config.Mode.Exec
	cmd, err := gophClient.Command(execCmd, os.Args[1:]...)
	if err != nil {
		log.Fatalf("创建命令失败: %v", err)
	}

	// 同步
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令并获取错误返回值
	err = cmd.Run()
	if err != nil {
		// 命令执行失败时返回错误码
		os.Exit(1)
	}
}
