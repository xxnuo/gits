package main

import (
	"flag"
	"log"
	"strings"

	"github.com/melbahja/goph"
)

var (
	configPath string
	config     *Config
)

func init() {
	flag.StringVar(&configPath, "c", "./config.toml", "配置文件路径, 默认使用当前目录下的 config.toml 文件")
}

func main() {
	flag.Parse()

	cmd := flag.Args()[0]

	println(cmd)

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("无法加载配置文件 %s: %v", configPath, err)
	}

	// 创建SSH连接
	client, err := goph.NewConn(&goph.Config{
		User:     config.SSH.User,
		Addr:     config.SSH.IP,
		Port:     config.SSH.Port,
		Auth:     goph.Password(config.SSH.Password),
		Callback: VerifyHost,
	})

	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	// 如果命令不是以git开头，自动添加git前缀
	if !strings.HasPrefix(cmd, "git") {
		cmd = "git " + cmd
	}
}
