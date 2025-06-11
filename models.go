package main

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/melbahja/goph"
)

// Config 结构体定义配置文件格式
type Config struct {
	SSH  SSHConfig  `toml:"ssh"`
	Repo RepoConfig `toml:"repo"`
}

// SSHConfig 定义SSH连接配置
type SSHConfig struct {
	Host string `toml:"host"`
	// override ssh config file settings
	Hostname       string `toml:"hostname"`
	User           string `toml:"user"`
	Port           uint   `toml:"port"`
	Password       string `toml:"password"`
	IdentityFile   string `toml:"identityfile"`
	Passphrase     string `toml:"passphrase"`
	ControlMaster  string `toml:"controlmaster"`
	ControlPath    string `toml:"controlpath"`
	ControlPersist string `toml:"controlpersist"`
	// internal settings
	HasAgent bool
}

// RepoConfig 定义仓库配置
type RepoConfig struct {
	Path string `toml:"path"`
}

// LoadConfig 从配置文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 规范化配置文件路径
	configPath, err := NormalizePath(configPath)
	if err != nil {
		return nil, errors.New("无法解析配置文件路径: " + err.Error())
	}

	// 解析配置文件
	var config Config
	if _, err = toml.DecodeFile(configPath, &config); err != nil {
		return nil, errors.New("读取配置文件失败: " + err.Error())
	}

	// 处理路径
	if err = NormalizeConfigPaths(&config); err != nil {
		return nil, err
	}

	// 验证配置
	if err = ValidateSSHConfig(&config.SSH); err != nil {
		return nil, err
	}

	// 处理内部设置
	config.SSH.HasAgent = goph.HasAgent()

	// 验证仓库配置
	if err = ValidateRepoConfig(&config.Repo); err != nil {
		return nil, err
	}

	return &config, nil
}

// NormalizeConfigPaths 处理配置中的所有路径
func NormalizeConfigPaths(config *Config) error {
	var err error

	// 处理 SSH 密钥路径
	if config.SSH.IdentityFile != "" {
		config.SSH.IdentityFile, err = NormalizePath(config.SSH.IdentityFile)
		if err != nil {
			return errors.New("无法解析密钥路径: " + err.Error())
		}
	}

	return nil
}

// ValidateSSHConfig 验证SSH配置的有效性
func ValidateSSHConfig(sshConfig *SSHConfig) error {
	if sshConfig.Host == "" {
		if sshConfig.Hostname == "" {
			return errors.New("SSH 服务器名称或地址不能同时为空")
		}
	} else {
		if sshConfig.IdentityFile == "" && sshConfig.Password == "" {
			return errors.New("SSH 密钥或密码不能同时为空")
		}
	}

	return nil
}

// ValidateRepoConfig 验证仓库配置的有效性
func ValidateRepoConfig(repoConfig *RepoConfig) error {
	if repoConfig.Path == "" {
		return errors.New("仓库路径不能为空")
	}

	return nil
}
