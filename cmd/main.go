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

	conf := configs.SetConfigFromEnv()

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

	// service.GetSaveNewCategories(ctx, s)
	// service.GetSaveNewIngredients(ctx, s)
	// service.GetSaveNewSubIngredients(ctx, s)
	// service.GetSaveNewRecepty(ctx, s)
	service.GetSaveFileRecept(ctx, s)

	return nil
}
