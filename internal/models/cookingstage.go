package models

type CookingStage struct {
	ID          int    `db:"id"`
	IDRecept    int    `db:"recept_id"`
	Number      string `db:"number"`
	Description string `db:"description"`
}

func (CookingStage) GetType() string {
	return "CookingStage"
}

func (c CookingStage) GetArgsInsert() []any {
	return []any{c.IDRecept, c.Number, c.Description}
}
