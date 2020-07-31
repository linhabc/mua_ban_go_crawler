package main

// run this function to crawl just mua ban o to category
// func crawlMuaBanOTo() {
// 	users := NewUsers()
// 	res := getHTMLPage("https://muaban.net/o-to-toan-quoc-l0-c4?cp=1")
// 	err := users.getAllUserInformation(res)
// 	checkError(err)
// 	users.TotalPages++

// 	for i := 2; i <= 200; i++ {

// 		users.TotalPages++
// 		nextPageLink := users.getNexURL(res)

// 		println(nextPageLink)

// 		res = getHTMLPage(nextPageLink)

// 		err := users.getAllUserInformation(res)
// 		checkError(err)
// 	}

// 	userJSON, err := json.Marshal(users) // convert User sang JSON
// 	checkError(err)
// 	err = ioutil.WriteFile("output_o_to.json", userJSON, 0644) // Ghi dữ liệu vào file JSON
// 	checkError(err)
// }
