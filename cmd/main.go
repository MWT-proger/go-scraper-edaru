package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Category struct {
	Slug       string
	Name       string
	Href       string
	ParentSlug string
}

func main() {
	// fmt.Println(getCategoryList())
	// fmt.Println(getReceptyList())
	fmt.Println(getRecepty())
}

func getCategoryList() []Category {
	listCategory := []Category{}

	c := colly.NewCollector(
		colly.AllowedDomains("eda.ru"),
	)
	c.SetProxy("94.228.194.18:41890")

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
			category := Category{
				Slug:       ParentSlug,
				Name:       strings.ReplaceAll(name, number, ""),
				Href:       href,
				ParentSlug: "",
			}
			listCategory = append(listCategory, category)
			// fmt.Println(category)
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

			category := Category{
				Slug:       slug,
				Name:       strings.ReplaceAll(name, number, ""),
				Href:       href,
				ParentSlug: ParentSlug,
			}

			// fmt.Println(category)
			listCategory = append(listCategory, category)
		})

	})
	c.Visit("https://eda.ru")

	return listCategory
}

func getReceptyList() int {
	urlCategory := "https://eda.ru/recepty/gribnoi-bulyon"

	count := 0
	allCount := 0
	first := true

	c := colly.NewCollector(
		colly.AllowedDomains("eda.ru"),
	)
	// c.SetProxy("94.228.194.18:41890")

	c.OnHTML(".emotion-1jdotsv", func(h *colly.HTMLElement) {
		if first {
			fmt.Println(h.Text)
			allCountString := h.Text
			listStr := []string{
				"Найдено ",
				" рецепта",
				" рецептов",
				" рецепт",
			}
			for _, v := range listStr {
				allCountString = strings.ReplaceAll(allCountString, v, "")
			}
			allCount, _ = strconv.Atoi(allCountString)

			first = false
		}
	})

	c.OnHTML(".emotion-1eugp2w", func(h *colly.HTMLElement) {
		categoryLinks := h.ChildAttrs("a", "href")
		if categoryLinks == nil {
			return
		}

		href := categoryLinks[0]
		fmt.Println(href)
		count++
	})
	c.Visit(urlCategory)
	if count != allCount {
		fmt.Println("На первой странице не все рецепты")
		for i := 2; count != allCount; i++ {
			url := fmt.Sprintf("%s?page=%s", urlCategory, strconv.Itoa(i))
			fmt.Println("Парсинг страницы № ", i)
			c.Visit(url)
		}
	}

	fmt.Println("Выведено рецептов: ", count)

	return count
}

type Recept struct {
	ID             int
	Name           string
	CookingTime    string
	Description    string
	NumberServings string
	ImageSrc       string
}

//	type Ingredient struct {
//		Name string
//	}
type IngredientRecept struct {
	IDRecept   int
	Quantity   string
	Ingredient string
}
type CookingStage struct {
	IDRecept    int
	Number      string
	Description string
}

func getRecepty() Recept {
	urlRecepty := "https://eda.ru/recepty/zavtraki/sirniki-iz-tvoroga-18506"
	s := strings.Split(urlRecepty, "-")
	id, _ := strconv.Atoi(s[len(s)-1])
	recept := Recept{ID: id}
	listIngredientRecepts := []IngredientRecept{}
	listCookingStage := []CookingStage{}

	c := colly.NewCollector(
		colly.AllowedDomains("eda.ru"),
	)

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
		ingredientRecept := IngredientRecept{
			IDRecept:   recept.ID,
			Quantity:   h.ChildText(".emotion-bsdd3p"),
			Ingredient: h.ChildText("span[itemprop=recipeIngredient]"),
		}
		listIngredientRecepts = append(listIngredientRecepts, ingredientRecept)

	})

	c.OnHTML("div .emotion-1lxj5xg", func(h *colly.HTMLElement) {

		cookingStage := CookingStage{
			IDRecept:    recept.ID,
			Number:      h.ChildText(".emotion-xhemb9"),
			Description: h.ChildText(".emotion-1dvddtv span"),
		}
		listCookingStage = append(listCookingStage, cookingStage)

	})

	c.Visit(urlRecepty)

	fmt.Println(listIngredientRecepts)
	fmt.Println(listCookingStage)

	return recept
}
