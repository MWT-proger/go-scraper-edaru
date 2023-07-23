package scraper

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/MWT-proger/go-scraper-edaru/internal/logger"
)

func (s *EdaRu) DowloadFile(hrefImage string, dirName string, fileName string) string {
	var (
		alowedContentType = map[string]string{"image/jpeg": "jpg", "image/png": "png"}
		conf              = configs.GetConfig()
		pathFile          = ""
		c                 = s.initColly()
	)

	c.OnResponse(func(r *colly.Response) {

		extension, ok := alowedContentType[r.Headers.Get("Content-Type")]

		if !ok {
			logger.Log.Error("Ссылка на картинку не имеет разрешенного Content-Type")
			return
		}

		pathFile = fmt.Sprintf("%s/%s/%s.%s", conf.BasePathDir, dirName, fileName, extension)

		if err := os.WriteFile(pathFile, r.Body, 0644); err != nil {
			pathFile = ""
			logger.Log.Error(err.Error())
			return
		}
	})

	c.Visit(hrefImage)

	return pathFile
}
