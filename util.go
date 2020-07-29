package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getHTMLPage(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		println("ERROR GET")
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		println("ERORR RES STATUS")
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil
	}
	return doc
}

func (users *Users) getNexURL(doc *goquery.Document) string {
	nextPageLink, _ := doc.Find("#next-link").Attr("href")

	print("NEXTPAGE: ")
	println(nextPageLink)
	// Trường hợp không có url
	if nextPageLink == "" {
		println("End of Category")
		return ""
	}

	return nextPageLink
}

func (users *Users) getAllUserInformation(doc *goquery.Document) error {
	doc.Find("a.list-item__link").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")
		go users.getUserInformation(userLink) // create goroutines
	})
	return nil
}

func (users *Users) getUserInformation(url string) {
	res := getHTMLPage(url)
	if res == nil {
		return
	}

	userName := res.Find(".user-info__fullname").Text()
	title := res.Find(".title").Text()
	price := res.Find(".price-container__value").Text()
	phoneNum, _ := res.Find("span[mobile]").Attr("mobile")

	user := User{
		URL:         url,
		UserName:    userName,
		Title:       title,
		Price:       price,
		PhoneNumber: phoneNum,
	}

	users.TotalUsers++
	users.List = append(users.List, user)
}

func checkError(err error) {
	if err != nil {
		print("Error: ")
		log.Println(err)
	}
}
