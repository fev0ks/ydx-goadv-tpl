package model

type User struct {
	Name string
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Current  int64
	Withdraw int64
}
