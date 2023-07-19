package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
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

	scr := scraper.EdaRu{Domen: "eda.ru"}
	// fmt.Println(scr.GetCategoryList())
	fmt.Println(scr.GetReceptyList("https://eda.ru/recepty/gribnoi-bulyon"))
	// fmt.Println(scr.GetRecepty("https://eda.ru/recepty/zavtraki/sirniki-iz-tvoroga-18506"))

	return nil
}
