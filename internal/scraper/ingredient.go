package scraper

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

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
				Description: sql.NullString{String: description, Valid: true},
				Href:        sql.NullString{String: href, Valid: true},
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
				Description: sql.NullString{String: description, Valid: true},
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
