package config

import "github.com/xeonds/phi-plug-go/lib"

type Config struct {
	// 服务器配置
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	// 数据库配置
	lib.DatabaseConfig `yaml:"database"`
	// 数据目录配置
	Data struct {
		Difficulty string `yaml:"difficulty"`
		Info       string `yaml:"info"`
	} `yaml:"data"`
}
