package scraper

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
)

type EdaRu struct {
	Domen string
}

func (s *EdaRu) initColly() *colly.Collector {

	c := colly.NewCollector(colly.AllowedDomains(s.Domen))
	c.Limit(&colly.LimitRule{
		DomainGlob: s.Domen,
		Delay:      1 * time.Second,
	})
	return c
}

func (s *EdaRu) GetCategoryList() []*models.Category {
	logger.Log.Debug("Парсинг категорий рецептов ...")
	listCategory := []*models.Category{}

	c := colly.NewCollector(
		colly.AllowedDomains(s.Domen),
	)

	c.OnHTML(".emotion-18mh8uc .emotion-w5dos9", func(e *colly.HTMLElement) {
		ParentSlug := ""
		e.ForEach(".emotion-w5dos9", func(_ int, h *colly.HTMLElement) {
			categoryLinks := h.ChildAttrs("a", "href")

			if categoryLinks == nil {
				return
			}

			href := categoryLinks[0]
			ParentSlug = strings.ReplaceAll(href, "/recepty/", "")
			name := h.ChildText("a h3")
			number := h.ChildText("a h3 span")
			category := models.Category{
				Slug: ParentSlug,
				Name: strings.ReplaceAll(name, number, ""),
				Href: href,
			}
			listCategory = append(listCategory, &category)
		})

		e.ForEach(".emotion-8asrz1", func(_ int, h *colly.HTMLElement) {
			categoryLinks := h.ChildAttrs("a", "href")
			if categoryLinks == nil {
				return
			}

			href := categoryLinks[0]
			slug := strings.ReplaceAll(href, "/recepty/", "")
			name := h.ChildText("a span")
			number := h.ChildText("a span span")

			category := models.Category{
				Slug:       slug,
				Name:       strings.ReplaceAll(name, number, ""),
				Href:       href,
				ParentSlug: sql.NullString{String: ParentSlug, Valid: true},
			}

			listCategory = append(listCategory, &category)
		})

	})
	c.Visit("https://" + s.Domen)

	logger.Log.Info("Категории получены", zap.Int("количество", len(listCategory)))
	return listCategory
}

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
			Quantity:   h.ChildText(".emotion-bsdd3p"),
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

// Парсит и выводит список ингредиентов основных
func (s *EdaRu) GetIngredientList() ([]*models.Ingredient, error) {
	logger.Log.Debug("Scraper: Парсинг ингредиентов ...")

	var (
		err            error
		listIngredient = []*models.Ingredient{}
		baseURL        = "https://" + s.Domen + "/wiki/"
		countLocal     = 0
		c              = colly.NewCollector(colly.AllowedDomains(s.Domen))
	)

	c.Limit(&colly.LimitRule{
		DomainGlob: s.Domen,
		Delay:      3 * time.Second,
	})

	var subC = c.Clone()
	subC.OnHTML(".emotion-17kxgoe", func(e *colly.HTMLElement) {
		// Получение списка игредиентов по категории

		var (
			name        = e.ChildText(".emotion-wxopay a h2")
			href        = e.ChildAttr(".emotion-wxopay a", "href")
			description = e.ChildText(".emotion-aus3ft")
			ss          = strings.Split(href, "-")
			id, _       = strconv.Atoi(ss[len(ss)-1])
		)

		if name != "" {
			countLocal++

			listIngredient = append(listIngredient, &models.Ingredient{
				ID:          id,
				Name:        name,
				Description: description,
				Href:        href,
				UpdatedAt:   time.Now(),
			})
		}

	})

	c.OnHTML(".emotion-mvfezh", func(e *colly.HTMLElement) {
		// Получение категорий ингредиентов и переход на неё
		var (
			attrs    = e.ChildAttrs(".emotion-1aag8k0 .emotion-1gwkmfw a", "href")
			count    = e.ChildText("a .emotion-a0cecn")
			countInt int
			listStr  = []string{
				" ингредиентов",
				" ингредиента",
				" ингредиент",
				" ",
			}
		)

		if attrs == nil {
			return
		}
		for _, v := range listStr {
			count = strings.ReplaceAll(count, v, "")
		}

		if count == "" {
			return
		}

		countInt, err = strconv.Atoi(count)

		if err != nil || countInt == 0 {
			return
		}
		countLocal = 0
		subC.Visit(baseURL + attrs[0])

		logger.Log.Info(
			"Scraper: Получение ингредиентов...",
			zap.Int("Получено", countLocal),
			zap.Int("Всего", countInt),
			zap.String("Категория", attrs[0]),
		)

		lastCount := 0

		if countLocal != countInt {

			for i := 2; lastCount != countLocal; i++ {

				lastCount = countLocal

				url := fmt.Sprintf("%s?page=%s", baseURL+attrs[0], strconv.Itoa(i))

				subC.Visit(url)

				logger.Log.Info(
					"Scraper: Получение ингредиентов...",
					zap.Int("Получено", countLocal),
					zap.Int("Всего", countInt),
					zap.String("Категория", attrs[0]),
				)
			}
		}

	})

	c.Visit(baseURL + "ingredienty")

	logger.Log.Info(
		"Scraper: Всего плучено Игредиентов",
		zap.Int("количество", len(listIngredient)))

	return listIngredient, nil
}

// Парсит и выводит список дочерних ингредиентов
func (s *EdaRu) GetSubIngredientList(urlParentIngredient string, parentID sql.NullInt64) ([]*models.Ingredient, error) {
	logger.Log.Debug("Scraper: Парсинг  дочерних ингредиентов ...")

	var (
		// err            error
		listIngredient = []*models.Ingredient{}
		baseURL        = "https://" + s.Domen
		countLocal     = 0
		c              = s.initColly()
	)

	c.OnHTML(".emotion-h0bpot .emotion-h0bpot", func(e *colly.HTMLElement) {
		// Получение категорий ингредиентов и переход на неё
		var (
			attrs       = e.ChildAttr("a", "href")
			name        = e.ChildText(".emotion-1cyq6dp")
			description = e.ChildText(".emotion-ketj7d span")
			ss          = strings.Split(attrs, "/")
			id, err     = strconv.Atoi(ss[len(ss)-1])
		)

		if err != nil {
			return
		}

		if name != "" {
			countLocal++

			listIngredient = append(listIngredient, &models.Ingredient{
				ID:          id,
				Name:        name,
				Description: description,
				UpdatedAt:   time.Now(),
				ParentId:    parentID,
			})

		}

	})

	c.Visit(baseURL + urlParentIngredient)

	logger.Log.Info(
		"Scraper: Всего плучено дочерних игредиентов",
		zap.Int("количество", len(listIngredient)))

	return listIngredient, nil
}
