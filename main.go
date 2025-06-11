package main

import (
	"log"
	"os"
	"strings"

	"github.com/melbahja/goph"
)

var (
	configPath string
	config     *Config
)

func main() {
	cmd := strings.Join(os.Args[1:], " ")
	if cmd == "" {
		return
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("无法加载配置文件 %s: %v", configPath, err)
	}

	hostKeyCallback, err := goph.DefaultKnownHosts()
	if err != nil {
		log.Fatalf("获取默认 known_hosts 文件失败: %v", err)
	}

	auth, err := AuthHost(config.SSH.HasAgent, config.SSH.Password, config.SSH.IdentityFile, config.SSH.Passphrase)
	if err != nil {
		log.Fatalf("认证失败: %v", err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     config.SSH.User,
		Addr:     config.SSH.Hostname,
		Auth:     auth,
		Callback: hostKeyCallback,
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
