package service

import (
	"context"
	"database/sql"

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
	categoryStorager.Insert(ctx, categories)

	return nil
}

func GetSaveNewSubIngredients(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr                = scraper.EdaRu{Domen: "eda.ru"}
		ingredientStorager = ingredientstorage.New(storage)
	)

	parentIngredients, err := ingredientStorager.GetByParameters(ctx, "SELECT * FROM content.ingredient WHERE parent_id is null", map[string]interface{}{"parent": ""})

	if err != nil {
		return err
	}

	for _, v := range parentIngredients {
		ingredients, err := scr.GetSubIngredientList(v.Href, sql.NullInt64{Int64: int64(v.ID), Valid: true})

		if err != nil {
			return err
		}
		ingredientStorager.Insert(ctx, ingredients)
	}

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
