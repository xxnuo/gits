package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
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

	// load config
	config, err = LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config file %s: %v", configPath, err)
	}

	gophConfig, err = GetHostConfig(config.SSH.Host)
	if err != nil {
		log.Fatalf("Failed to get host config: %v", err)
	}

	gophClient, err = goph.NewConn(gophConfig)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer gophClient.Close()

	execCmd := "cd " + config.Repo.Path + " && " + config.Mode.Exec

	session, err := gophClient.Client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
		os.Exit(1)
	}
	defer session.Close()
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// 检查当前终端是否为交互式终端
	isTty := term.IsTerminal(int(os.Stdin.Fd()))
	var oldState *term.State
	var fd int
	if isTty {
		// 对于交互式终端，创建一个支持PTY的SSH会话
		fd = int(os.Stdin.Fd())
		var err error
		oldState, err = term.MakeRaw(fd)
		if err != nil {
			log.Fatalf("Failed to make raw: %v", err)
		}
		defer term.Restore(fd, oldState)

		width, height, err := term.GetSize(fd)
		if err != nil {
			term.Restore(fd, oldState)
			log.Fatalf("Failed to get size: %v", err)
		}

		err = session.RequestPty("xterm-256color", height, width, ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		})
		if err != nil {
			term.Restore(fd, oldState)
			log.Fatalf("Failed to request PTY: %v", err)
		}
	}
	err = session.Run(execCmd + " " + strings.Join(os.Args[1:], " "))
	if err != nil && isTty {
		term.Restore(fd, oldState)
		os.Exit(1)
	} else if err != nil {
		os.Exit(1)
	}
}
