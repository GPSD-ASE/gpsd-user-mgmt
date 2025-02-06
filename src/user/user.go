package user

type User struct {
	Name  string `json:"name"`
	DevID string `json:"devID"`
}

var users = []User{
	User{"abc", "123"},
	User{"qwe", "121"},
	User{"zxc", "134"},
}
