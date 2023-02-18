package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	NewStatus        = OrderStatus("NEW")
	ProcessingStatus = OrderStatus("PROCESSING")
	InvalidStatus    = OrderStatus("INVALID")
	ProcessedStatus  = OrderStatus("PROCESSED")
)

type OrderStatus string

type Order struct {
	Number     int `db:"order_id"`
	Status     OrderStatus
	Accrual    int
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Number     string
		Status     OrderStatus
		Accrual    float32
		UploadedAt string `json:"uploaded_at"`
	}{
		Number:     fmt.Sprintf("%d", o.Number),
		Status:     o.Status,
		Accrual:    float32(o.Accrual) / 100,
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}

type AccrualOrder struct {
	Order   int
	Status  OrderStatus
	Accrual int
}

func (ao *AccrualOrder) UnmarshalJSON(data []byte) error {
	var accrualOrderIn struct {
		Order   string
		Status  OrderStatus
		Accrual *float32
	}
	if err := json.Unmarshal(data, &accrualOrderIn); err != nil {
		return err
	}
	orderNumber, err := strconv.Atoi(accrualOrderIn.Order)
	if err != nil {
		return err
	}
	ao.Order = orderNumber
	ao.Status = accrualOrderIn.Status
	if accrualOrderIn.Accrual != nil {
		ao.Accrual = int(*accrualOrderIn.Accrual * 100)
	} else {
		ao.Accrual = 0
	}
	return nil
}
func (ao *AccrualOrder) ToOrder() *Order {
	return &Order{
		Number:  ao.Order,
		Status:  ao.Status,
		Accrual: ao.Accrual,
	}
}

type UserOrder struct {
	UserId int
	Order  *Order
}
