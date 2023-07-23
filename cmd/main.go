package main

import (
	"context"
	"log"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/service"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

// initProject() иницилизирует все необходимые переменный проекта
func initProject(ctx context.Context) (*storage.PgStorage, error) {

	var (
		configInit = configs.InitConfig()
	)

	parseFlags(configInit)

	configs.SetConfigFromEnv()

	conf, err := configs.ValidateConfig()

	if err != nil {
		return nil, err
	}

	if err := logger.Initialize(conf.LogLevel); err != nil {
		return nil, err
	}

	s := &storage.PgStorage{}

	if err := s.Init(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func main() {
	log.Println("Запуск проекта")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		cancel()
		panic(err)
	}
}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера
func run(ctx context.Context) error {
	s, err := initProject(ctx)

	if err != nil {
		return err
	}

	defer s.Close()

	scenario, err := configs.GetScenario()

	if err != nil {
		return err
	}

	if v, ok := scenario["ingredients"]; ok && v == 1 {
		log.Println("Парсинг и загрузка Ингредиентов ВКЛ")
		service.GetSaveNewIngredients(ctx, s)
	}
	if v, ok := scenario["sub_ingredients"]; ok && v == 1 {
		log.Println("Парсинг и загрузка Дочерних ингредиентов ВКЛ")
		service.GetSaveNewSubIngredients(ctx, s)
	}
	if v, ok := scenario["categories"]; ok && v == 1 {
		log.Println("Парсинг и загрузка Категорий ВКЛ")
		service.GetSaveNewCategories(ctx, s)
	}
	if v, ok := scenario["recipes"]; ok && v == 1 {
		log.Println("Парсинг и загрузка Рецептов ВКЛ")
		service.GetSaveNewRecepty(ctx, s)
	}
	if v, ok := scenario["file_recipes"]; ok && v == 1 {
		log.Println("Парсинг и загрузка Файлов рецептов ВКЛ")
		service.GetSaveFileRecept(ctx, s)
	}

	return nil
}
