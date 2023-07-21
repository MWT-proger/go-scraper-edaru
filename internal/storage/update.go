package storage

import (
	"context"
	"database/sql"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UpdatePgStorage[E models.BaseModeler] struct {
	*PgStorage
}

type Updateer[E models.BaseModeler] interface {
	Update(ctx context.Context, updateQuery string, tx *sql.Tx, obj E, args map[string]interface{})
}

func NewUpdatePgStorage[E models.BaseModeler](baseStorage *PgStorage) *UpdatePgStorage[E] {
	return &UpdatePgStorage[E]{baseStorage}
}

// (s *UpdatePgStorage[E]) Update(obj E) Это базовый метод
// для добавления объектов в хранилище
func (s *UpdatePgStorage[E]) Update(ctx context.Context, updateQuery string, txBig *sqlx.Tx, obj E, args map[string]interface{}) error {
	logger.Log.Info("Обновление строки в хранилище данных...", zap.String("таблица", obj.GetType()))
	var bigTx bool
	var tx *sqlx.Tx
	var err error

	if txBig != nil {
		bigTx = true
		tx = txBig
	}

	if !bigTx {

		tx, err = s.db.BeginTxx(ctx, nil)

		if err != nil {
			return err
		}

		defer tx.Rollback()
	}

	stmt, err := tx.PrepareNamedContext(ctx, updateQuery)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args)

	if err != nil {

		logger.Log.Error(err.Error())
		return err
	}

	if !bigTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	logger.Log.Info("Строка обновлена")

	return nil

}
