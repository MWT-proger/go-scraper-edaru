package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"go.uber.org/zap"
)

type InsertPgStorage[E models.BaseModeler] struct {
	*PgStorage
	insertQuery string
}

type Inserter[E models.BaseModeler] interface {
	Insert(ctx context.Context, objs []E)
}

func NewInsertPgStorage[E models.BaseModeler](baseStorage *PgStorage, insertQuery string) *InsertPgStorage[E] {
	return &InsertPgStorage[E]{baseStorage, insertQuery}
}

// (s *InsertPgStorage[E]) Insert(obj E) Это базовый метод
// для добавления объектов в хранилище
func (s *InsertPgStorage[E]) Insert(ctx context.Context, objs []E) error {
	logger.Log.Info("Добавление в хранилище данных...", zap.Int("количество", len(objs)))

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, s.insertQuery)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	count := 0
	for _, v := range objs {

		res, err := stmt.ExecContext(ctx, v.GetArgsInsert()...)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Log.Debug("Обект уже существует в хранилище")
				continue
			}
			logger.Log.Error(err.Error())
			return err
		}

		if r, _ := res.RowsAffected(); r == 1 {
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Log.Info("Добавлены в хранилище новые данные", zap.Int("количество", count))

	return nil

}
