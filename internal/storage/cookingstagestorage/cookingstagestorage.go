package cookingstagestorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type CookingStageStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.CookingStage]
	*storage.GetByParametersPgStorage[*models.CookingStage]
}

type CookingStageStorager interface {
	storage.Inserter[models.CookingStage]
}

func New(baseStorage *storage.PgStorage) *CookingStageStorage {
	insertQuery := "INSERT INTO content.cooking_stage (recept_id, number, description) VALUES($1,$2,$3) RETURNING (id)"
	insertRepo := storage.NewInsertPgStorage[*models.CookingStage](baseStorage, insertQuery)
	geterRepo := storage.NewGetByParametersPgStorage[*models.CookingStage](baseStorage)

	return &CookingStageStorage{baseStorage, insertRepo, geterRepo}
}
