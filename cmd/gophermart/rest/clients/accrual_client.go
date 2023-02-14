package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts/rest"
	"github.com/pkg/errors"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	ordersAPI = "/api/orders/{number}"
)

type AccrualClient interface {
	GetOrderStatus(ctx context.Context, orderNumber int) (*model.AccrualOrder, error)
}

type accrualClient struct {
	client *resty.Client
}

func NewAccrualClient(client *resty.Client) AccrualClient {
	return &accrualClient{client}
}

func (ac accrualClient) GetOrderStatus(ctx context.Context, orderNumber int) (*model.AccrualOrder, error) {
	resp, err := ac.client.R().
		SetHeader(rest.ContentType, rest.TextPlain).
		SetPathParams(map[string]string{
			"number": fmt.Sprintf("%d", orderNumber),
		}).
		SetContext(ctx).
		Get(ordersAPI)
	if err != nil {
		return nil, err
	}
	return ac.parseOrderStatusResponse(resp)
}

func (ac accrualClient) parseOrderStatusResponse(resp *resty.Response) (*model.AccrualOrder, error) {
	if resp.StatusCode()/100 != 2 {
		respBody := resp.Body()
		return nil, errors.Errorf("Accrual OrderStatus response status is not successfull: '%s', body: '%s'", resp.Status(), strings.TrimSpace(string(respBody)))
	}
	log.Printf("Process Order status is %s", resp.Status())
	accrualOrder := &model.AccrualOrder{}
	err := json.Unmarshal(resp.Body(), &accrualOrder)
	if err != nil {
		return nil, err
	}
	log.Printf("Processed Order is '%v'", accrualOrder)
	return accrualOrder, nil
}
