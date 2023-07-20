package models

type Ingredient struct {
	Name        string
	Description string
}

type IngredientRecept struct {
	IDRecept   int
	Quantity   string
	Ingredient string
}
