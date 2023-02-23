package clients

import (
	"github.com/go-resty/resty/v2"
	"time"
)

func CreateClient(baseURL string) *resty.Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	return client
}
