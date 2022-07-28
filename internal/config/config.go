package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"pinger/logger"
	"sync"
)

type Config struct {
	IsDebug *bool         `yaml:"is_debug" env-required:"true"`
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var instance *Config
var once sync.Once

// GetConfig - Функция чтения данных из конфига
func GetConfig() *Config {
	once.Do(func() {
		logger := logger.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
