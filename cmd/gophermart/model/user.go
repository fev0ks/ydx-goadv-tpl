package model

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"login"`
	Password []byte `json:"password"`
}

type Balance struct {
	Current  int64
	Withdraw int64
}
