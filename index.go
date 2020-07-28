package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func crawlFromCategory(category Category) {
	users := NewUsers()
	res := getHTMLPage(category.URL)

	//handle error
	if res == nil {
		return
	}
	err := users.getAllUserInformation(res)
	checkError(err)
	users.TotalPages++

	for i := 2; i <= 200; i++ {
		users.TotalPages++
		nextPageLink := users.getNexURL(res)

		if nextPageLink == "" {
			break
		}

		res = getHTMLPage(nextPageLink)

		//handle error
		if res == nil {
			break
		}

		err := users.getAllUserInformation(res)
		checkError(err)
	}

	// convert User sang JSON
	userJSON, err := json.Marshal(users)
	checkError(err)

	// Ghi dữ liệu vào file JSON
	err = ioutil.WriteFile(category.Title+".json", userJSON, 0644)
	checkError(err)
}

func main() {

	file, _ := ioutil.ReadFile("categories.json")

	data := Categories{}

	_ = json.Unmarshal([]byte(file), &data)

	for i := 4; i < len(data.List); i++ {
		fmt.Println("Title: ", data.List[i].Title)
		fmt.Println("URL: ", data.List[i].URL)

		crawlFromCategory(data.List[i])
	}

}
