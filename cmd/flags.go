package main

import (
	"flag"

	"github.com/MWT-proger/go-scraper-edaru/configs"
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags(conf *configs.Config) {
	flag.StringVar(&conf.DatabaseDSN, "d", conf.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&conf.LogLevel, "l", "info", "уровень логирования")
	flag.StringVar(&conf.BasePathDir, "p", conf.BasePathDir, "путь каталога для сохранения файлов")
	flag.Parse()
}
