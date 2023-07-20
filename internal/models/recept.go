package models

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

type CookingStage struct {
	IDRecept    int
	Number      string
	Description string
}
