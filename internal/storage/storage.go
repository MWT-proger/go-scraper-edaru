package storage

import (
	"context"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PgStorage struct {
	db *sqlx.DB
}

// Init(ctx context.Context) error Инициализирует соединение
// и возвращает ошибку в случае не удачи
// Вызывается при запуске программы
func (s *PgStorage) Init(ctx context.Context) error {
	logger.Log.Info("Хранилище: Подключение...")

	var (
		conf    = configs.GetConfig()
		db, err = sqlx.Open("pgx", conf.DatabaseDSN)
	)

	if err != nil {
		return err
	}

	s.db = db

	if err := s.ping(); err != nil {
		return err
	}

	if err := s.migration(); err != nil {
		return err
	}
	logger.Log.Info("Хранилище: Соединение установленно")

	return nil

}

func (s *PgStorage) GetDB() *sqlx.DB {
	return s.db
}

// Migration() проверяет новые миграции и при неообходимости добавляет в БД
// и возвращает ошибку в случае не удачи
// Вызывается при запуске программы
func (s *PgStorage) migration() error {
	logger.Log.Info("Хранилище: Проверка и обновление миграций ...")

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}

// ping() error прверяет соединение и возвращает ошибку в случае не удачи
func (s *PgStorage) ping() error {
	logger.Log.Info("Хранилище: Проверка соединения ...")
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

// Close() error закрывает соединение и возвращает ошибку в случае не удачи
// Вызывается при завершение программы
func (s *PgStorage) Close() error {
	logger.Log.Info("Хранилище: Закрытие соединения...")

	if err := s.db.Close(); err != nil {
		return err
	}
	logger.Log.Info("Хранилище: Соединение успешно закрыто")

	return nil
}
