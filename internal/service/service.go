package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/categorystorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/cookingstagestorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/ingredientreceptstorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/ingredientstorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/receptcategorystorage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/receptstorage"
)

func GetSaveNewCategories(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr              = scraper.EdaRu{Domen: "eda.ru"}
		categories       = scr.GetCategoryList()
		categoryStorager = categorystorage.New(storage)
	)
	if err := categoryStorager.Insert(ctx, nil, categories); err != nil {
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
		ingredients, err := scr.GetSubIngredientList(v.Href.String, sql.NullInt64{Int64: int64(v.ID), Valid: true})

		if err != nil {
			return err
		}

		if err := ingredientStorager.Insert(ctx, nil, ingredients); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	return nil
}

func GetSaveNewRecepty(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr                      = scraper.EdaRu{Domen: "eda.ru"}
		categorystorager         = categorystorage.New(storage)
		receptystorager          = receptstorage.New(storage)
		ingredientreceptstorager = ingredientreceptstorage.New(storage)
		receptcategoryer         = receptcategorystorage.New(storage)
		cookingstagestorager     = cookingstagestorage.New(storage)
		ingredientStorager       = ingredientstorage.New(storage)
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

		tx, err := storage.GetDB().BeginTx(ctx, nil)

		if err != nil {
			return err
		}

		defer tx.Rollback()

		if err := receptystorager.Insert(ctx, tx, recepties); err != nil {
			logger.Log.Error(err.Error())
			return err
		}

		for _, v := range recepties {
			for _, ing := range v.Ingredients {
				ingredient, err := ingredientStorager.GetByParameters(ctx, "SELECT * FROM content.ingredient WHERE name=:ingredient_name LIMIT 1", map[string]interface{}{"ingredient_name": ing.Ingredient})

				if err != nil {
					logger.Log.Error(err.Error())
					return err
				}
				if len(ingredient) == 0 {
					time.Sleep(1 * time.Second)
					newIngredient := models.Ingredient{ID: int(time.Now().Unix()), Name: ing.Ingredient}

					ingredientStorager.Insert(ctx, tx, []*models.Ingredient{&newIngredient})

					ing.IngredientID = newIngredient.ID
				} else {
					ing.IngredientID = ingredient[0].ID
				}

			}

			ingredientreceptstorager.Insert(ctx, tx, v.Ingredients)

			receptcategoryer.Insert(ctx, tx, []*models.ReceptCategory{{ReceptID: v.ID, CategorySlug: v.CategorySlug}})
			cookingstagestorager.Insert(ctx, tx, v.CookingStages)
		}
		if err := tx.Commit(); err != nil {
			return err
		}

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

	if err := ingredientStorager.Insert(ctx, nil, ingredients); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func GetSaveFileRecept(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr             = scraper.EdaRu{Domen: "eda.ru"}
		receptystorager = receptstorage.New(storage)
	)

	recepties, err := receptystorager.GetByParameters(ctx,
		"SELECT id, image_src FROM content.recept WHERE image IS NULL AND image_src IS NOT NULL ORDER BY id;",
		map[string]interface{}{})

	if err != nil {
		return err
	}
	for _, v := range recepties {

		if v.ImageSrc == "" {
			continue
		}

		str := scr.DowloadFile(v.ImageSrc, "recepties", strconv.Itoa(v.ID))

		if str == "" {
			return errors.New("file не загружен")
		}

		v.Image = str

		if err := receptystorager.Update(ctx,
			"UPDATE content.recept SET image=:path WHERE id=:id_recept", nil, v,
			map[string]interface{}{"path": str, "id_recept": v.ID}); err != nil {
			return err
		}

	}

	return nil
}
