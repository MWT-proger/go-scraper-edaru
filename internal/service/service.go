package service

import (
	"context"

	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/categorystorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/ingredientstorage"
)

// fmt.Println(scr.GetReceptyList("https://eda.ru/recepty/gribnoi-bulyon"))

func GetSaveNewCategories(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr              = scraper.EdaRu{Domen: "eda.ru"}
		categories       = scr.GetCategoryList()
		categoryStorager = categorystorage.New(storage)
	)
	// categories := []*models.Category{}
	// categories = append(categories, &models.Category{Slug: "123", Name: "123", Href: "234"})
	categoryStorager.Insert(ctx, categories)

	return nil
}

func GetSaveNewIngredients(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr                = scraper.EdaRu{Domen: "eda.ru"}
		ingredients, err   = scr.GetIngredientList()
		ingredientStorager = ingredientstorage.New(storage)
	)

	if err != nil {
		return err
	}
	// categories := []*models.Ingredient{}
	// categories = append(categories, &models.Ingredient{ID: 23, Name: "123", Href: "234", UpdatedAt: time.Now()})

	ingredientStorager.Insert(ctx, ingredients)

	return nil
}

func GetSaveNewRecepty(ctx context.Context, storage *storage.PgStorage) error {

	var (
		scr = scraper.EdaRu{Domen: "eda.ru"}
	)
	scr.GetRecepty("https://eda.ru/recepty/vypechka-deserty/brauni-brownie-20955")

	return nil
}
