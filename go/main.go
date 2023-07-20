package main 

import (
	"github.com/gocolly/colly/v2"
)

const URL = "https://cd.mfa.gov.tr/mission/mission-list?clickedId=3"
const STORE_FILE = "./store.json"
func main() {
	c := colly.NewCollector(
		colly.AllowedDomains()
	)
	file, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	store := NewStore(STORE_FILE, AccountInfo{})
	c.OnHTML("li[data-griddercontent]", func(e *colly.HTMLElement) {
		div := e.Id("MissionsList")
		country_id := e.Attr("data-griddercontent").Split("=")[1]		
		country_name := e.ChildText("span")
	})

	c.OnHTML("li[data-griddercontent]>span", func(e *colly.HTMLElement) {
		e.ChildText()
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(URL)

}