package model

import "time"

type Withdraw struct {
	Order       int
	Sum         int64
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRequest struct {
	Order int
	Sum   int64
}
