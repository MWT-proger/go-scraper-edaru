package main

import (
	"fmt"

	"github.com/MWT-proger/go-scraper-edaru/internal/scraper"
)

func main() {
	s := scraper.EdaRu{Domen: "eda.ru"}
	fmt.Println(s.GetCategoryList())
	fmt.Println(s.GetReceptyList("https://eda.ru/recepty/gribnoi-bulyon"))
	fmt.Println(s.GetRecepty("https://eda.ru/recepty/zavtraki/sirniki-iz-tvoroga-18506"))
}
