package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

func crawlFromCategory(category Category) {
	// open leveldb connection
	db := createOrOpenDb("./db/" + category.Title)
	defer db.Close()

	users := NewUsers()
	res := getHTMLPage(category.URL)

	//handle error
	if res == nil {
		return
	}

	err := users.getAllUserInformation(res, db)
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

		err := users.getAllUserInformation(res, db)
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

func crawlAllFromCategories(categories Categories) {
	var wg sync.WaitGroup

	jobs := make(chan Category, 100)

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(w, jobs)
	}

	for i := 0; i < len(categories.List); i++ {
		jobs <- categories.List[i]
	}

	close(jobs)

	wg.Wait()
}

func worker(id int, jobs <-chan Category) {
	for j := range jobs {
		fmt.Println("worker: ", id, "processing job: ", j)
		crawlFromCategory(j)
	}
}

func main() {
	// get categories from json file
	file, _ := ioutil.ReadFile("categories.json")

	data := Categories{}

	_ = json.Unmarshal([]byte(file), &data)

	// schedule to run each 3 hour
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for true {
			crawlAllFromCategories(data)
			time.Sleep(3 * time.Hour)
		}
		wg.Done()
	}()

	wg.Wait()
}
