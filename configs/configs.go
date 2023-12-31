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
		BasePathDir:  "data",
		FileScenario: "",
	}
	return &newConfig
}

// GetConfig() выводит не импортируемую переменную newConfig
func GetConfig() Config {
	return newConfig
}

// SetConfigFromEnv() Присваевает полям значения из ENV
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

// GetScenario() выводит сценапий парсинга
func GetScenario() (map[string]int, error) {
	data := make(map[string]int)

	if newConfig.FileScenario == "" {
		data["ingredients"] = 1
		data["sub_ingredients"] = 1
		data["categories"] = 1
		data["recipes"] = 1
		data["file_recipes"] = 1
		return data, nil
	}
	f, err := os.ReadFile(newConfig.FileScenario)

	if err != nil {
		return nil, errors.New("сценарий не найден")
	}

	if err := json.Unmarshal([]byte(f), &data); err != nil {
		return nil, errors.New("не верный формат сценария")
	}

	return data, nil

}

// GetScenario() проверяет обязательные переменные
func ValidateConfig() (Config, error) {

	if newConfig.DatabaseDSN == "" {
		return Config{}, errors.New("необходимо указать параметры подключения к БД")

	}

	return newConfig, nil
}
