package configs

import (
	"os"
)

type Config struct {
	LogLevel     string
	DatabaseDSN  string `env:"DATABASE_DSN"`
	DomenScraper string `env:"DOMEN_SCRAPER"`
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		LogLevel:     "info",
		DatabaseDSN:  "",
		DomenScraper: "eda.ru",
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
	return newConfig
}
