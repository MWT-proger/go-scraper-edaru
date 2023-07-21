package scraper

import (
	"fmt"
	"os"

	"github.com/MWT-proger/go-scraper-edaru/configs"
	"github.com/gocolly/colly"
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
			return
		}

		pathFile = fmt.Sprintf("%s/%s/%s.%s", conf.BasePathDir, dirName, fileName, extension)

		if err := os.WriteFile(pathFile, r.Body, 0644); err != nil {
			pathFile = ""
			fmt.Println(err)
			return
		}
	})

	c.Visit(hrefImage)

	return pathFile
}
