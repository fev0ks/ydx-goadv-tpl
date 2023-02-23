package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Withdraw struct {
	WithdrawID  int `json:"-" db:"withdraw_id"`
	Order       int
	Sum         int
	ProcessedAt time.Time `json:"processed_at"`
}

func (wd *Withdraw) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Order       string
		Sum         float32
		ProcessedAt string `json:"processed_at"`
	}{
		Order:       fmt.Sprintf("%d", wd.Order),
		Sum:         float32(wd.Sum) / 100,
		ProcessedAt: wd.ProcessedAt.Format(time.RFC3339),
	})
}

type WithdrawRequest struct {
	Order int
	Sum   int
}

func (wd *WithdrawRequest) UnmarshalJSON(data []byte) error {
	var withdrawRequestIn struct {
		Order string
		Sum   float32
	}
	if err := json.Unmarshal(data, &withdrawRequestIn); err != nil {
		return err
	}
	orderNumber, err := strconv.Atoi(withdrawRequestIn.Order)
	if err != nil {
		return err
	}
	wd.Order = orderNumber
	wd.Sum = int(withdrawRequestIn.Sum * 100)
	return nil
}
