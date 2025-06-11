package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	gitCmd := strings.Join(os.Args[1:], " ")
	if gitCmd == "" {
		return
	}

	config, err = LoadConfig(configPath)
	if err != nil {
		log.Fatalf("无法加载配置文件 %s: %v", configPath, err)
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

	// 如果命令不是以git开头，自动添加git前缀
	if !strings.HasPrefix(gitCmd, "git") {
		gitCmd = "git " + gitCmd
	}

	gitCmd = "cd " + config.Repo.Path + " && " + gitCmd

	cmd, err := gophClient.Command(gitCmd)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("执行命令失败: %v", err)
	}

	fmt.Println(cmd.Stdout)

}
