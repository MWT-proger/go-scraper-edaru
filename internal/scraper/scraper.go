package scraper

import (
	"time"

	"github.com/gocolly/colly"
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
