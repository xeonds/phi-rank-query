package config

import "github.com/xeonds/phi-plug-go/lib"

type Config struct {
	// 服务器配置
	Server struct {
		Port               string `yaml:"port"`
		InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
		LogFile            string `yaml:"logFile"`
	}
	// 数据库配置
	lib.DatabaseConfig
	// 数据目录配置
	Data struct {
		Difficulty string `yaml:"difficulty"`
		Info       string `yaml:"info"`
		Version    string `yaml:"version"`
	}
}
