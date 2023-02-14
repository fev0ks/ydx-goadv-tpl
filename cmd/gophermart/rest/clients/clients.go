package clients

import (
	"github.com/go-resty/resty/v2"
	"time"
)

func CreateClient(baseUrl string) *resty.Client {
	client := resty.New().
		SetBaseURL(baseUrl).
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(2 * time.Second)
	return client
}
