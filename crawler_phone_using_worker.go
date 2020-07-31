package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
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
	dt := time.Now()
	err = ioutil.WriteFile("./output/"+category.Title+dt.String()+".json", userJSON, 0644)
	checkError(err)
}

func worker(id int, jobs <-chan Category) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		crawlFromCategory(j)
	}
}

func crawlAllFromCategories(data Categories) {
	var wg sync.WaitGroup

	jobs := make(chan Category, 100)

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(w, jobs)
	}

	for i := 0; i < len(data.List); i++ {
		fmt.Println("Title: ", data.List[i].Title)
		fmt.Println("URL: ", data.List[i].URL)

		jobs <- data.List[i]
	}

	close(jobs)

	wg.Wait()
}

func checkExist() {

}

func main() {

	// get json data
	file, _ := ioutil.ReadFile("categories.json")

	data := Categories{}

	_ = json.Unmarshal([]byte(file), &data)

	db := make([]*leveldb.DB, len(data.List))

	// open levelDb
	for i := 0; i < len(data.List); i++ {
		println(data.List[i].Title)
		db[i] = createOrOpenDb("./db/" + data.List[i].Title)
		defer db[i].Close()
	}

	// schedule to run each 6 hour
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for true {
			crawlAllFromCategories(data)
			time.Sleep(6 * time.Hour)
		}
		wg.Done()
	}()

	wg.Wait()
}
