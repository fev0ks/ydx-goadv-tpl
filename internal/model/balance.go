package model

import "encoding/json"

type Balance struct {
	userID   int `json:"-" db:"user_id"`
	Current  int64
	Withdraw int64
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Current  float32 `json:"current"`
		Withdraw float32 `json:"withdrawn"`
	}{
		Current:  float32(b.Current) / 100,
		Withdraw: float32(b.Withdraw) / 100,
	})
}
