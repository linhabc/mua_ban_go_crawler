package main

// Category is a
type Category struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Categories is a
type Categories struct {
	Total int        `json:"total"`
	List  []Category `json:"categories"`
}

func newCategories() *Categories {
	return &Categories{}
}

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
