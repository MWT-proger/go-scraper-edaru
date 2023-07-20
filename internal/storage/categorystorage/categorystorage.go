package categorystorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type CategoryStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.Category]
}

type CategoryStorager interface {
	storage.Inserter[models.Category]
}

func New(baseStorage *storage.PgStorage) *CategoryStorage {
	insertQuery := "INSERT INTO content.category (slug, name, href, parent_slug) VALUES($1,$2,$3,$4) ON CONFLICT (slug) DO NOTHING RETURNING (slug)"
	insertRepo := storage.NewInsertPgStorage[*models.Category](baseStorage, insertQuery)

	return &CategoryStorage{baseStorage, insertRepo}
}
