package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

// Config 结构体定义配置文件格式
type Config struct {
	SSH  SSHConfig  `toml:"ssh"`
	Repo RepoConfig `toml:"repo"`
	Mode ModeConfig `toml:"mode"`
}

// SSHConfig 定义SSH连接配置
type SSHConfig struct {
	Host string `toml:"host"`
}

// RepoConfig 定义仓库配置
type RepoConfig struct {
	Path string `toml:"path"`
}

// ModeConfig 定义模式配置
type ModeConfig struct {
	Exec string `toml:"exec"`
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

	// 验证配置
	if err = ValidateSSHConfig(&config.SSH); err != nil {
		return nil, err
	}

	// 验证仓库配置
	if err = ValidateRepoConfig(&config.Repo); err != nil {
		return nil, err
	}

	// 验证模式配置
	if err = ValidateModeConfig(&config.Mode); err != nil {
		return nil, err
	}

	return &config, nil
}

// ValidateSSHConfig 验证SSH配置的有效性
func ValidateSSHConfig(sshConfig *SSHConfig) error {
	if sshConfig.Host == "" {
		return errors.New("SSH 服务器名称不能为空")
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

// ValidateModeConfig 验证模式配置的有效性
func ValidateModeConfig(modeConfig *ModeConfig) error {
	if modeConfig.Exec == "" {
		modeConfig.Exec = "git"
	}

	return nil
}
