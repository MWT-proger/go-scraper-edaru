package service

import (
	"context"

	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage"
	"github.com/MWT-proger/go-scraper-edaru/internal/storage/categorystorage"
)

// fmt.Println(scr.GetReceptyList("https://eda.ru/recepty/gribnoi-bulyon"))
// fmt.Println(scr.GetRecepty("https://eda.ru/recepty/zavtraki/sirniki-iz-tvoroga-18506"))

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

func GetSaveNewIngredient(ctx context.Context, storage *storage.PgStorage) error {
	var (
		scr = scraper.EdaRu{Domen: "eda.ru"}
	)
	scr.GetIngredientList()

	return nil
}
