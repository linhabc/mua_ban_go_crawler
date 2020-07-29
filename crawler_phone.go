package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
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
	err = ioutil.WriteFile("./output/"+category.Title+".json", userJSON, 0644)
	checkError(err)
}

// func worker(jobs <-chan int) {
// 	for j := range jobs {

// 	}
// }

func main() {
	// jobs := make(chan int)

	var wg sync.WaitGroup
	file, _ := ioutil.ReadFile("categories.json")

	data := Categories{}

	_ = json.Unmarshal([]byte(file), &data)

	for i := 1; i < len(data.List); i++ {
		wg.Add(1)
		fmt.Println("Title: ", data.List[i].Title)
		fmt.Println("URL: ", data.List[i].URL)

		go func(i int) {
			crawlFromCategory(data.List[i])
			wg.Done()
		}(i)
	}

	wg.Wait()
}
