package domain

type User struct {
    Account string `json:"account"`
	Password string `json:"password"`
	Name string `json:"name"`
	Email string `json:"email"`
}