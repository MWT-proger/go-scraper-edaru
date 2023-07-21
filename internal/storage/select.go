package storage

import (
	"context"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"go.uber.org/zap"
)

type GetByParametersPgStorage[E models.BaseModeler] struct {
	*PgStorage
}

type GetByParameterser[E models.BaseModeler] interface {
	GetByParameters(ctx context.Context, selectQuery string, args map[string]interface{}) ([]E, error)
}

func NewGetByParametersPgStorage[E models.BaseModeler](baseStorage *PgStorage) *GetByParametersPgStorage[E] {
	return &GetByParametersPgStorage[E]{baseStorage}
}

// (s *GetByParametersPgStorage[E]) GetByParameters(obj E) Это базовый метод
// для добавления объектов в хранилище
func (s *GetByParametersPgStorage[E]) GetByParameters(ctx context.Context, selectQuery string, args map[string]interface{}) ([]E, error) {
	logger.Log.Info("Получение из хранилища данных...")

	stmt, err := s.db.PrepareNamedContext(ctx, selectQuery)
	list := []E{}

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &list, args); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	logger.Log.Info("Получены объекты", zap.Int("количество", len(list)))

	return list, nil

}
