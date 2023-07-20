package models

type IngredientRecept struct {
	IDRecept   int    `db:"recept_id"`
	Quantity   string `db:"quantity"`
	Ingredient string `db:"ingredient_id"`
}

func (IngredientRecept) GetType() string {
	return "IngredientRecept"
}

func (c IngredientRecept) GetArgsInsert() []any {
	return []any{c.IDRecept, c.Quantity, c.Ingredient}
}
