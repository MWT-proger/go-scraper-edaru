package receptcategorystorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type ReceptCategoryStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.ReceptCategory]
	*storage.GetByParametersPgStorage[*models.ReceptCategory]
}

type ReceptCategoryStorager interface {
	storage.Inserter[models.ReceptCategory]
}

func New(baseStorage *storage.PgStorage) *ReceptCategoryStorage {
	insertQuery := "INSERT INTO content.recept_category (recept_id, category_slug) VALUES($1,$2) ON CONFLICT (recept_id, category_slug) DO NOTHING RETURNING (recept_id, category_slug)"
	insertRepo := storage.NewInsertPgStorage[*models.ReceptCategory](baseStorage, insertQuery)
	geterRepo := storage.NewGetByParametersPgStorage[*models.ReceptCategory](baseStorage)

	return &ReceptCategoryStorage{baseStorage, insertRepo, geterRepo}
}
