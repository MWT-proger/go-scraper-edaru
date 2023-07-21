package receptstorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type ReceptStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.Recept]
	*storage.GetByParametersPgStorage[*models.Recept]
	*storage.UpdatePgStorage[*models.Recept]
}

type ReceptStorager interface {
	storage.Inserter[models.Recept]
}

func New(baseStorage *storage.PgStorage) *ReceptStorage {
	insertQuery := "INSERT INTO content.recept (id, name, cooking_time, description, number_servings, image_src) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO NOTHING RETURNING (id)"
	insertRepo := storage.NewInsertPgStorage[*models.Recept](baseStorage, insertQuery)
	geterRepo := storage.NewGetByParametersPgStorage[*models.Recept](baseStorage)
	updaterRepo := storage.NewUpdatePgStorage[*models.Recept](baseStorage)

	return &ReceptStorage{baseStorage, insertRepo, geterRepo, updaterRepo}
}
