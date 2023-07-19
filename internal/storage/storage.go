package storage

import (
	"context"
	"database/sql"
	"embed"

	"github.com/gocolly/colly/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgStorage struct {
	storage.Storage
	db *sql.DB
}

func (s *PgStorage) Init(ctx context.Context) error {
	conf := configs.GetConfig()
	logger.Log.Info("Подключение к БД ...")
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return err
	}
	s.db = db

	if err := s.Ping(); err != nil {
		return err
	}

	if err := s.Migration(); err != nil {
		return err
	}
	logger.Log.Info("Соединение с БД установленно")

	return nil

}

// Migration() проверяет новые миграции и при неообходимости добавляет в БД
func (s *PgStorage) Migration() error {
	logger.Log.Info("Проверка и обновление миграций ...")
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Ping() error {
	logger.Log.Info("Проверка соединения ...")
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Close() error {
	logger.Log.Info("Закрытие соединения с БД ...")

	if err := s.db.Close(); err != nil {
		return err
	}
	logger.Log.Info("Соединение успешно закрыто")

	return nil
}
