package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/syndtr/goleveldb/leveldb"
)

func getHTMLPage(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		println("ERROR GET")
		return nil
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

func (users *Users) getAllUserInformation(doc *goquery.Document, db *leveldb.DB) error {
	doc.Find("a.list-item__link").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")

		// create goroutines
		go users.getUserInformation(userLink, db)
	})
	return nil
}

func (users *Users) getUserInformation(url string, db *leveldb.DB) {
	res := getHTMLPage(url)
	if res == nil {
		return
	}

	// check if id-url is exit in db or not
	checkExist := getData(db, url)
	if checkExist != "" {
		println("Exist: " + url)
		return
	}

	phoneNum, _ := res.Find("span[mobile]").Attr("mobile")

	user := User{
		Id:          url,
		PhoneNumber: phoneNum,
	}

	putData(db, url, phoneNum)

	users.TotalUsers++
	users.List = append(users.List, user)
}

func checkError(err error) {
	if err != nil {
		print("Error: ")
		log.Println(err)
	}
}
