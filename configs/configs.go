package configs

import (
	"os"
)

type Config struct {
	LogLevel     string
	DatabaseDSN  string `env:"DATABASE_DSN"`
	DomenScraper string `env:"DOMEN_SCRAPER"`
	BasePathDir  string `env:"BASE_PATH_DIR"`
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		LogLevel:     "info",
		DatabaseDSN:  "",
		DomenScraper: "eda.ru",
		BasePathDir:  "./../data",
	}
	return &newConfig
}

// GetConfig() выводит не импортируемую переменную newConfig
func GetConfig() Config {
	return newConfig
}

// SetConfigFromEnv() Прсваевает полям значения из ENV
// Вызывается один раз при старте проекта
func SetConfigFromEnv() Config {
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		newConfig.LogLevel = envLogLevel
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		newConfig.DatabaseDSN = envDatabaseDSN
	}
	if envDomenScraper := os.Getenv("DOMEN_SCRAPER"); envDomenScraper != "" {
		newConfig.DomenScraper = envDomenScraper
	}
	if envBasePathDir := os.Getenv("BASE_PATH_DIR"); envBasePathDir != "" {
		newConfig.BasePathDir = envBasePathDir
	}
	return newConfig
}
