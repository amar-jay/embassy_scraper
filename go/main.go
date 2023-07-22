package main

import (
	. "scraper/store"
)

const URL = "https://cd.mfa.gov.tr/mission/"
const STORE_FILE = "storage.json"

func main() {

	// var c = colly.NewCollector(
	// 	colly.AllowedDomains(),
	// )
	var store = NewStore[AccountInfo](STORE_FILE)
	defer store.Close()

	// scrape(c, store)
	// file_creation(c, store)
	// info_scrape(c, store)
	// country_names(store)

	if err := store.Save(); err != nil {
		panic(err)
	}
	if err := store.Flush(); err != nil {
		panic(err)
	}
}
