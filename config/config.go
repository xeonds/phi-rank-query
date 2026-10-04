package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port               string `yaml:"port" mapstructure:"port"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify" mapstructure:"insecureSkipVerify"`
	LogFile            string `yaml:"logFile" mapstructure:"logFile"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type" mapstructure:"type"`
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       string `yaml:"db" mapstructure:"db"`
	Migrate  bool   `yaml:"migrate" mapstructure:"migrate"`
}

type DataConfig struct {
	Difficulty string `yaml:"difficulty" mapstructure:"difficulty"`
	Info       string `yaml:"info" mapstructure:"info"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"databaseconfig" mapstructure:"databaseconfig"`
	Data     DataConfig     `yaml:"data" mapstructure:"data"`
}

// Load reads config.yaml from the working directory. If it is missing, a
// template is written and the program exits so it can be filled in.
func Load() *Config {
	if _, err := os.Stat("config.yaml"); err != nil {
		data, _ := yaml.Marshal(new(Config))
		_ = os.WriteFile("config.yaml", data, 0o644)
		log.Fatal("config file not found, a template has been created")
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config file read failed: ", err)
	}
	cfg := new(Config)
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal("config file parse failed: ", err)
	}
	return cfg
}
