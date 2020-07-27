package main

import (
	"github.com/PuerkitoBio/goquery"
)

// Category is a
type Category struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Categories is a
type Categories struct {
	Total int        `json:"total"`
	List  []Category `json:"categories"`
}

func newCategories() *Categories {
	return &Categories{}
}

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

// func main() {
// 	categories := newCategories()
// 	res := getHTMLPage("https://muaban.net/")
// 	categories.getAllCategories(res)

// 	userJSON, err := json.Marshal(categories)
// 	checkError(err)
// 	err = ioutil.WriteFile("categories.json", userJSON, 0644) // Ghi dữ liệu vào file JSON
// 	checkError(err)
// }
