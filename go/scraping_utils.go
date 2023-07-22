package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func scrape(c *colly.Collector, store Store[AccountInfo]) {

	c.OnHTML("li[data-griddercontent]", func(e *colly.HTMLElement) {
		country_id := strings.Split(e.Attr("data-griddercontent"), "=")[1]
		country_name := e.ChildText("span.dp-name")
		country_name = strings.Trim(country_name, " ")
		store.Set(country_id, AccountInfo{
			Name: country_name,
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(URL + "mission-list?clickedId=3")

}

func file_creation(c *colly.Collector, store Store[AccountInfo]) {

	data, err := store.GetAll()
	if err != nil {
		panic(err)
	}
	for i, ai := range data {

		perm := fs.FileMode(0777)
		// join "details" directory with country name
		fmt.Printf("%+v\n", ai)
		file := fmt.Sprintf("html/%s_%02d.html", ai.Name, i)
		f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
		if err != nil {
			panic(err)
		}
		log.Printf("created file %s", file)
		f.Close()
	}

	log.Print("files created")
	// fmt.Printf("%d: %+v\n", i, ai)

}

func info_scrape(c *colly.Collector, store Store[AccountInfo]) {
	data, err := store.GetMap()
	if err != nil {
		panic(err)
	}
	c.OnHTML("#embassy>address", func(e *colly.HTMLElement) {
		text := e.Text
		country_id := strings.Split(e.Request.URL.String(), "=")[1]

		sb := strings.Builder{}
		spl := strings.Split(text, "\n")
		ai, _ := store.Get(country_id)
		for _, line := range spl {
			line = strings.Trim(line, " ")
			if len(line) == 0 {
				continue
			}
			// fetch email if line contains email
			if strings.Contains(line, "@") {
				email := strings.Trim(line, ": ")
				ai.Email = append(ai.Email, email)
			}
			sb.WriteString(line + "\n")
		}

		file := fmt.Sprintf("details/%s.txt", ai.Name)
		perm := fs.FileMode(0777)
		outliersFile, _ := os.OpenFile("outliers.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, perm)
		f, _ := os.OpenFile(file, os.O_CREATE|os.O_RDWR, perm)
		_, err := f.Write([]byte(sb.String()))
		if err != nil {
			panic(err)
		}
		f.Close()

		text = sb.String()
		spl = strings.Split(text, "\n")
		if len(spl) < 7 {
			outliersFile.WriteString(fmt.Sprintf("%s %s\n", country_id, ai.Name))
			return
		}

		ai.Address = Address(text)
		ai.Ambassador = Ambassador(text)
		ai.Phone = Phone(text)
		ai.Email = email(text)
		store.Set(country_id, ai)
	})

	for i := range data {

		// join "details" directory with country name
		err = c.Visit(URL + "MissionDetail?CountryId=" + i)
		if err != nil {
			panic(err)
		}
	}

}

func country_names(store Store[AccountInfo]) {
	// print country names

	data, err := store.GetAll()
	if err != nil {
		panic(err)
	}

	for _, ai := range data {
		fmt.Println(ai.Name)
	}

}
