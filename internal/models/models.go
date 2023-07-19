package models

type Category struct {
	Slug       string
	Name       string
	Href       string
	ParentSlug string
}

type Recept struct {
	ID             int
	Name           string
	CookingTime    string
	Description    string
	NumberServings string
	ImageSrc       string
	Ingredients    []IngredientRecept
	CookingStages  []CookingStage
}

type Ingredient struct {
	Name string
}

type IngredientRecept struct {
	IDRecept   int
	Quantity   string
	Ingredient string
}
type CookingStage struct {
	IDRecept    int
	Number      string
	Description string
}
