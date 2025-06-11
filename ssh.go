package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	knownHostsPath, err := goph.DefaultKnownHostsPath()
	if err != nil {
		return fmt.Errorf("获取known_hosts文件路径失败: %v", err)
	}

	// 检查主机是否在known_hosts文件中
	hostFound, err := goph.CheckKnownHost(host, remote, key, knownHostsPath)

	// 如果主机存在但验证失败，可能是主机密钥已更改，返回错误
	if hostFound && err != nil {
		log.Printf("警告: 主机 %s 的密钥验证失败: %v", host, err)
		return err
	}

	// 如果主机存在且验证成功，直接返回
	if hostFound && err == nil {
		return nil
	}

	// 主机不存在，询问用户是否添加
	fmt.Printf("主机 %s 不在 %s 文件中。\n", host, knownHostsPath)
	fmt.Printf("主机指纹: %s\n", ssh.FingerprintSHA256(key))
	fmt.Printf("是否添加到 %s? [y/n]: ", knownHostsPath)

	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		// 添加主机到known_hosts文件
		return goph.AddKnownHost(host, remote, key, knownHostsPath)
	}

	return fmt.Errorf("主机验证被用户拒绝")
}

func AuthHost(hasAgent bool, password string, identityFile string, passphrase string) (goph.Auth, error) {
	if hasAgent {
		return goph.UseAgent()
	}
	if identityFile != "" {
		return goph.Key(identityFile, passphrase)
	}
	if password != "" {
		return goph.Password(password), nil
	}

	return nil, nil
}
