package model

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Balance struct {
	Current  int64
	Withdraw int64
}
