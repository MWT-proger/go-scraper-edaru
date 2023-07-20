package ingredientstorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type IngredientStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.Ingredient]
	*storage.GetByParametersPgStorage[*models.Ingredient]
}

type IngredientStorager interface {
	storage.Inserter[models.Ingredient]
}

func New(baseStorage *storage.PgStorage) *IngredientStorage {
	insertQuery := "INSERT INTO content.ingredient (id, name, description, href, parent_id, updated_at) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO NOTHING RETURNING (id)"
	insertRepo := storage.NewInsertPgStorage[*models.Ingredient](baseStorage, insertQuery)
	geterRepo := storage.NewGetByParametersPgStorage[*models.Ingredient](baseStorage)

	return &IngredientStorage{baseStorage, insertRepo, geterRepo}
}
