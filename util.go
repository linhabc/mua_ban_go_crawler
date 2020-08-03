package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

func (users *Users) getAllUserInformation(doc *goquery.Document, category string, f *os.File, db *leveldb.DB) error {
	doc.Find("a.list-item__link").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")

		// create goroutines
		go users.getUserInformation(userLink, category, f, db)
	})
	return nil
}

func (users *Users) getUserInformation(url string, category string, f *os.File, db *leveldb.DB) {
	res := getHTMLPage(url)
	if res == nil {
		return
	}

	splitResult := strings.Split(url, "-")
	id := splitResult[len(splitResult)-1]

	// check if id(url) is exit in db or not
	checkExist := getData(db, id)
	if checkExist != "" {
		println("Exist: " + id)
		return
	}

	phoneNum, _ := res.Find("span[mobile]").Attr("mobile")

	user := User{
		Id:          id,
		PhoneNumber: phoneNum,
	}

	putData(db, id, phoneNum)

	// convert User sang JSON
	userJSON, err := json.Marshal(user)

	checkError(err)
	io.WriteString(f, string(userJSON)+"\n")

	users.TotalUsers++
	users.List = append(users.List, user)
}

func checkError(err error) {
	if err != nil {
		print("Error: ")
		log.Println(err)
	}
}
