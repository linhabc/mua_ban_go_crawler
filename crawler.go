package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// User is a
type User struct {
	URL         string `json:"url"`
	UserName    string `json:"user_name"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	PhoneNumber string `json:"phone_number"`
}

// Users is a
type Users struct {
	TotalPages int    `json:"total_pages"`
	TotalUsers int    `json:"total_users"`
	List       []User `json:"users"`
}

// NewUsers is a
func NewUsers() *Users {
	return &Users{}
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func getHTMLPage(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
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

	// Trường hợp không có url
	if nextPageLink == "javascript:void();" {
		users.TotalPages = 1
		return ""
	}
	return nextPageLink
}

func (users *Users) getAllUserInformation(doc *goquery.Document) error {
	doc.Find("a.list-item__link").Each(func(i int, s *goquery.Selection) {
		userLink, _ := s.Attr("href")
		go users.getUserInformation(userLink) // create goroutines
		// println(userLink)
	})
	return nil
}

func (users *Users) getUserInformation(url string) {

	res := getHTMLPage(url)
	userName := res.Find(".user-info__fullname").Text()
	title := res.Find(".title").Text()
	price := res.Find(".price-container__value").Text()
	phoneNum, _ := res.Find("span[mobile]").Attr("mobile")
	// println("Name: " + userName)
	// println("title " + title)
	// println("price: " + price)
	// println("phoneNum: " + phoneNum)

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

func main() {
	users := NewUsers()
	res := getHTMLPage("https://muaban.net/o-to-toan-quoc-l0-c4?cp=1")
	err := users.getAllUserInformation(res)
	checkError(err)
	users.TotalPages++

	for i := 2; i <= 200; i++ {

		users.TotalPages++
		nextPageLink := users.getNexURL(res)

		println(nextPageLink)

		res = getHTMLPage(nextPageLink)

		// pageNum := strconv.Itoa(i)
		// res := getHTMLPage("https://muaban.net/o-to-toan-quoc-l0-c4?cp=" + pageNum)
		// println("https://muaban.net/o-to-toan-quoc-l0-c4?cp=" + pageNum)

		err := users.getAllUserInformation(res)
		checkError(err)
	}

	userJSON, err := json.Marshal(users) // convert User sang JSON
	checkError(err)
	err = ioutil.WriteFile("output_all_new.json", userJSON, 0644) // Ghi dữ liệu vào file JSON
	checkError(err)
}
