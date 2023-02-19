package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts/rest"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	ordersAPI = "/api/orders/{number}"
)

type AccrualClient interface {
	GetOrderStatus(ctx context.Context, orderID int) (*model.AccrualOrder, error)
}

type accrualClient struct {
	client *resty.Client
}

func NewAccrualClient(client *resty.Client) AccrualClient {
	return &accrualClient{client}
}

func (ac accrualClient) GetOrderStatus(ctx context.Context, orderID int) (*model.AccrualOrder, error) {
	resp, err := ac.client.R().
		SetHeader(rest.ContentType, rest.TextPlain).
		SetPathParams(map[string]string{
			"number": fmt.Sprintf("%d", orderID),
		}).
		SetContext(ctx).
		Get(ordersAPI)
	if err != nil {
		return nil, err
	}

	log.Printf("Process Order status is %s", resp.Status())
	if resp.StatusCode()/100 != 2 {
		if resp.StatusCode() == http.StatusTooManyRequests {
			retryAfterV := resp.Header().Get(rest.RetryAfter)
			retryAfter, err := strconv.Atoi(retryAfterV)
			if err != nil {
				log.Printf("failed to get RetryAfter value: %v", err)
				retryAfter = 1
			}
			log.Printf("Process Order is available after %d sec", retryAfter)
			return nil, &model.RetryError{
				Err:        errors.New("Too Many Requests"),
				RetryAfter: retryAfter,
			}
		} else {
			respBody := resp.Body()
			return nil, errors.Errorf("Accrual OrderStatus response status is not successfull: '%s', body: '%s'",
				resp.Status(), strings.TrimSpace(string(respBody)))
		}
	}
	accrualOrder := &model.AccrualOrder{}
	if resp.StatusCode() == http.StatusNoContent {
		log.Printf("Order is not presented in accrual system '%d'", orderID)
		accrualOrder = &model.AccrualOrder{
			Order:   orderID,
			Status:  model.NewStatus,
			Accrual: 0,
		}
	} else {
		err = json.Unmarshal(resp.Body(), &accrualOrder)
		if err != nil {
			return nil, err
		}
	}
	log.Printf("Processed Order is '%v'", accrualOrder)
	return accrualOrder, nil
}
