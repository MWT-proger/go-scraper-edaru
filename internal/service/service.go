package service

import (
	"context"
	"database/sql"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/categorystorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/ingredientstorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/receptstorage"
)

// fmt.Println(scr.GetReceptyList("https://eda.ru/recepty/gribnoi-bulyon"))

func GetSaveNewCategories(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr              = scraper.EdaRu{Domen: "eda.ru"}
		categories       = scr.GetCategoryList()
		categoryStorager = categorystorage.New(storage)
	)
	if err := categoryStorager.Insert(ctx, categories); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

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

func GetSaveNewRecepty(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr              = scraper.EdaRu{Domen: "eda.ru"}
		categorystorager = categorystorage.New(storage)
		receptystorager  = receptstorage.New(storage)
	)

	categories, err := categorystorager.GetByParameters(ctx, "SELECT * FROM content.category", map[string]interface{}{})

	if err != nil {
		return err
	}

	for _, v := range categories {
		recepties, err := scr.GetReceptyList(v.Href, v.Slug)

		if err != nil {
			return err
		}

		for _, v := range recepties {
			if err := scr.GetRecepty(v); err != nil {
				return err
			}
		}
		receptystorager.Insert(ctx, recepties)
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
