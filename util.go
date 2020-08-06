package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

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

func (users *Users) getAllUserInformation(doc *goquery.Document, category string, f *os.File, db *leveldb.DB) {
	var wg sync.WaitGroup
	doc.Find("a.list-item__link").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")
		wg.Add(1)
		go users.getUserInformation(userLink, category, &wg, f, db)
	})
	wg.Wait()
}

func (users *Users) getUserInformation(url string, category string, wg *sync.WaitGroup, f *os.File, db *leveldb.DB) {
	defer wg.Done()

	res := getHTMLPage(url)
	if res == nil {
		return
	}

	phoneNum, _ := res.Find("span[mobile]").Attr("mobile")

	if len(phoneNum) == 0 {
		return
	}

	splitResult := strings.Split(url, "-")
	id := splitResult[len(splitResult)-1]

	// check if id is exist in db or not
	checkExist := getData(db, id)
	if len(checkExist) != 0 {
		println("Exist: " + id)
		return
	} else {
		println("None_exist: " + id)
	}

	user := User{
		Id:          id,
		PhoneNumber: phoneNum,
	}

	_ = putData(db, id, phoneNum)

	// convert User to JSON
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
