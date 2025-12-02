package models

type LoginJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string
}
