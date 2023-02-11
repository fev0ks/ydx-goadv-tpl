package model

import "time"

const (
	NewStatus        = OrderStatus("NEW")
	ProcessingStatus = OrderStatus("PROCESSING")
	InvalidStatus    = OrderStatus("INVALID")
	ProcessedStatus  = OrderStatus("PROCESSED")
)

type OrderStatus string

type Order struct {
	Number     int
	UploadedAt time.Time `json:"uploaded_at"`
	Status     OrderStatus
}
