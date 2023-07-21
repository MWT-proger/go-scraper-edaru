package models

import "database/sql"

type IngredientRecept struct {
	IDRecept     int            `db:"recept_id"`
	Quantity     sql.NullString `db:"quantity"`
	IngredientID int            `db:"ingredient_id"`
	Ingredient   string
}

func (IngredientRecept) GetType() string {
	return "IngredientRecept"
}

func (c IngredientRecept) GetArgsInsert() []any {

	return []any{c.IDRecept, c.Quantity, c.IngredientID}
}
