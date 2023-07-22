package configs

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	LogLevel     string
	DatabaseDSN  string `env:"DATABASE_DSN"`
	DomenScraper string `env:"DOMEN_SCRAPER"`
	BasePathDir  string `env:"BASE_PATH_DIR"`
	FileScenario string `env:"FILE_SCENARIO"`
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
		FileScenario: "./../scenario.json",
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
	if envFileScenario := os.Getenv("FILE_SCENARIO"); envFileScenario != "" {
		newConfig.FileScenario = envFileScenario
	}
	return newConfig
}

func GetScenario() (map[string]interface{}, error) {
	var data map[string]interface{}
	f, err := os.ReadFile(newConfig.FileScenario)

	if err != nil {
		return nil, errors.New("сценарий не найден")
	}

	if err := json.Unmarshal([]byte(f), &data); err != nil {
		return nil, errors.New("не верный формат сценария")
	}

	return data, nil

}
