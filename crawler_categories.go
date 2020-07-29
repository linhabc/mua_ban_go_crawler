package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

func (categories *Categories) getAllCategories(doc *goquery.Document) {
	doc.Find(".category-item").Each(func(i int, s *goquery.Selection) {
		catLink, _ := s.Attr("href")
		catTitle, _ := s.Attr("title")

		category := Category{
			Title: catTitle,
			URL:   catLink,
		}

		categories.Total++
		categories.List = append(categories.List, category)
	})
}

// run this function to crawlAllCategories
func crawlAllCategories() {
	categories := newCategories()
	res := getHTMLPage("https://muaban.net/")
	categories.getAllCategories(res)

	userJSON, err := json.Marshal(categories)
	checkError(err)
	err = ioutil.WriteFile("./output/categories.json", userJSON, 0644) // Ghi dữ liệu vào file JSON
	checkError(err)
}
