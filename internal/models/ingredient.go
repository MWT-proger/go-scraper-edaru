package models

import "time"

type Ingredient struct {
	ID          int
	Name        string
	Description string
	Href        string
	ParentId    int
	UpdatedAt   time.Time
}

func (Ingredient) GetType() string {
	return "Ingredient"
}

func (c Ingredient) GetArgsInsert() []any {
	if c.ParentId != 0 {
		return []any{c.ID, c.Name, c.Description, c.Href, c.ParentId, c.UpdatedAt}
	} else {
		return []any{c.ID, c.Name, c.Description, c.Href, nil, c.UpdatedAt}
	}
}

type IngredientRecept struct {
	IDRecept   int
	Quantity   string
	Ingredient string
}
