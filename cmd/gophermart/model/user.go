package model

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	UserId   int    `json:"-" db:"user_id"`
	Username string `json:"login"`
	Password []byte `json:"password"`
}
