package ingredientreceptstorage

import (
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
)

type IngredientReceptStorage struct {
	*storage.PgStorage
	*storage.InsertPgStorage[*models.IngredientRecept]
	*storage.GetByParametersPgStorage[*models.IngredientRecept]
}

type IngredientReceptStorager interface {
	storage.Inserter[models.IngredientRecept]
}

func New(baseStorage *storage.PgStorage) *IngredientReceptStorage {
	insertQuery := "INSERT INTO content.ingredient_recept (recept_id, quantity, ingredient_id) VALUES($1,$2, (SELECT  id AS ingredient_id FROM content.ingredient WHERE name=$3)) " +
		"ON CONFLICT (recept_id, ingredient_id) DO NOTHING " +
		"RETURNING (recept_id, ingredient_id)"
	insertRepo := storage.NewInsertPgStorage[*models.IngredientRecept](baseStorage, insertQuery)
	geterRepo := storage.NewGetByParametersPgStorage[*models.IngredientRecept](baseStorage)

	return &IngredientReceptStorage{baseStorage, insertRepo, geterRepo}
}
