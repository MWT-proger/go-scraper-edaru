package models

type Recept struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	CookingTime    string `db:"cooking_time"`
	Description    string `db:"description"`
	NumberServings string `db:"number_servings"`
	ImageSrc       string `db:"image_src"`
	Href           string
	CategorySlug   string
	Ingredients    []*IngredientRecept
	CookingStages  []*CookingStage
}

func (Recept) GetType() string {
	return "Recept"
}

func (c Recept) GetArgsInsert() []any {
	return []any{c.ID, c.Name, c.CookingTime, c.Description, c.NumberServings, c.ImageSrc}
}
