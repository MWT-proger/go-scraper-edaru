package scraper

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

func (s *EdaRu) GetReceptyList(urlCategory string, slugCategory string) ([]*models.Recept, error) {
	logger.Log.Info(
		"Scraper: Получение списка ссылок на рецепты...",
	)
	var (
		count    = 0
		allCount = 0
		first    = true
		baseURL  = "https://" + s.Domen
		c        = s.initColly()
		list     = []*models.Recept{}
	)

	c.OnHTML(".emotion-1jdotsv", func(h *colly.HTMLElement) {
		if first {
			allCountString := h.Text
			listStr := []string{
				"Найдено ",
				"Найден ",
				"Найдены ",
				" рецепта",
				" рецептов",
				" рецепт",
			}
			for _, v := range listStr {
				allCountString = strings.ReplaceAll(allCountString, v, "")
			}
			allCount, _ = strconv.Atoi(allCountString)

			first = false
			logger.Log.Info(
				"Scraper: Всего рецептов по категории",
				zap.String("Категория", slugCategory),
				zap.Int("Всего", allCount),
			)
		}
	})

	c.OnHTML(".emotion-1eugp2w", func(h *colly.HTMLElement) {
		categoryLinks := h.ChildAttrs("a", "href")

		if categoryLinks == nil {
			return
		}

		list = append(list, &models.Recept{Href: categoryLinks[0], CategorySlug: slugCategory})

		count++
	})

	c.Visit(baseURL + urlCategory)

	lastCount := 0

	if count != allCount {
		fmt.Println("На первой странице не все рецепты")

		for i := 2; lastCount != count; i++ {
			lastCount = count

			url := fmt.Sprintf("%s?page=%s", baseURL+urlCategory, strconv.Itoa(i))

			logger.Log.Info(
				"Scraper: Получение рецептов...",
				zap.Int("Получено", count),
				zap.Int("Всего", allCount),
				zap.String("Категория", slugCategory),
			)
			fmt.Println("Парсинг страницы № ", i)

			c.Visit(url)
		}
	}
	logger.Log.Info(
		"Scraper: Получено рецептов",
		zap.Int("Получено", count),
		zap.Int("Всего", allCount),
		zap.String("Категория", slugCategory),
	)
	fmt.Println("Выведено рецептов: ", count)

	return list, nil
}

func (s *EdaRu) GetRecepty(recept *models.Recept) error {
	var (
		ss                    = strings.Split(recept.Href, "-")
		id, _                 = strconv.Atoi(ss[len(ss)-1])
		baseURL               = "https://" + s.Domen
		listIngredientRecepts = []*models.IngredientRecept{}
		listCookingStage      = []*models.CookingStage{}
		c                     = s.initColly()
	)
	recept.ID = id

	c.OnHTML("span[itemprop=resultPhoto]", func(h *colly.HTMLElement) {
		recept.ImageSrc = h.Attr("content")

	})
	c.OnHTML(".emotion-19rdt1j", func(h *colly.HTMLElement) {
		recept.Name = h.ChildText("h1")
		recept.CookingTime = h.ChildText(".emotion-my9yfq")
		recept.NumberServings = h.ChildText("span[itemprop=recipeYield]")

	})

	c.OnHTML(".emotion-aiknw3", func(h *colly.HTMLElement) {
		recept.Description = h.Text

	})

	c.OnHTML("div .emotion-7yevpr", func(h *colly.HTMLElement) {
		ingredientRecept := models.IngredientRecept{
			IDRecept:   recept.ID,
			Quantity:   sql.NullString{String: h.ChildText(".emotion-bsdd3p"), Valid: true},
			Ingredient: h.ChildText("span[itemprop=recipeIngredient]"),
		}
		listIngredientRecepts = append(listIngredientRecepts, &ingredientRecept)

	})

	c.OnHTML("div .emotion-1lxj5xg", func(h *colly.HTMLElement) {

		cookingStage := models.CookingStage{
			IDRecept:    recept.ID,
			Number:      h.ChildText(".emotion-xhemb9"),
			Description: h.ChildText(".emotion-1dvddtv span"),
		}
		listCookingStage = append(listCookingStage, &cookingStage)

	})

	c.Visit(baseURL + recept.Href)

	recept.CookingStages = listCookingStage
	recept.Ingredients = listIngredientRecepts
	fmt.Println(recept)
	return nil
}
