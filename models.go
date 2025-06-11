package main

import (
	"errors"

	"github.com/BurntSushi/toml"
)

// Config 结构体定义配置文件格式
type Config struct {
	SSH  SSHConfig  `toml:"ssh"`
	Repo RepoConfig `toml:"repo"`
}

// SSHConfig 定义SSH连接配置
type SSHConfig struct {
	// ssh connection settings
	IP             string `toml:"ip"`
	User           string `toml:"user"`
	Port           uint   `toml:"port"`
	Password       string `toml:"password"`
	Key            string `toml:"key"`
	ControlMaster  string `toml:"controlmaster"`
	ControlPath    string `toml:"controlpath"`
	ControlPersist string `toml:"controlpersist"`
	// ssh config file settings
	Path string `toml:"path"`
	Name string `toml:"name"`
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
	if config.SSH.Key != "" {
		config.SSH.Key, err = NormalizePath(config.SSH.Key)
		if err != nil {
			return errors.New("无法解析密钥路径: " + err.Error())
		}
	}

	// 处理 SSH 配置文件路径
	if config.SSH.Path != "" {
		config.SSH.Path, err = NormalizePath(config.SSH.Path)
		if err != nil {
			return errors.New("无法解析SSH配置文件路径: " + err.Error())
		}
	}

	return nil
}

// ValidateSSHConfig 验证SSH配置的有效性
func ValidateSSHConfig(sshConfig *SSHConfig) error {
	if sshConfig.Path != "" {
		// 使用 SSH 配置文件模式
		if sshConfig.Name == "" {
			return errors.New("SSH 服务器名称不能为空")
		}
	} else {
		// 使用直接配置模式
		if sshConfig.IP == "" {
			return errors.New("SSH IP 不能为空")
		}
		if sshConfig.User == "" {
			return errors.New("SSH 用户名不能为空")
		}
		if sshConfig.Key == "" && sshConfig.Password == "" {
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
