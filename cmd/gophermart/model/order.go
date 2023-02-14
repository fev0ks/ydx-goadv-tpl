package model

const (
	NewStatus        = OrderStatus("NEW")
	ProcessingStatus = OrderStatus("PROCESSING")
	InvalidStatus    = OrderStatus("INVALID")
	ProcessedStatus  = OrderStatus("PROCESSED")
)

type OrderStatus string

type Order struct {
	Number     string
	Status     OrderStatus
	Accrual    float32
	UploadedAt string `json:"uploaded_at"`
}

type AccrualOrder struct {
	Order   string
	Status  OrderStatus
	Accrual float32
}
