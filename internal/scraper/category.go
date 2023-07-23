package scraper

import (
	"database/sql"
	"strings"

	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
	"github.com/MWT-proger/go-scraper-edaru/internal/models"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

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
