package repository

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

type OrderRepository interface {
	InsertOrUpdateOrder(ctx context.Context, username string, order *model.Order) error
	UpdateOrderStatus(ctx context.Context, number int, newStatus model.OrderStatus) error
	IsOrderExist(ctx context.Context, orderNumber int) (bool, error)
	GetOrders(ctx context.Context, username string) ([]*model.Order, error)
}

type orderRepository struct {
	db DBProvider
}

func NewOrderRepository(db DBProvider) OrderRepository {
	return &orderRepository{db}
}

func (or orderRepository) InsertOrUpdateOrder(ctx context.Context, username string, order *model.Order) error {
	log.Printf("Persisting Order: %v", order)
	_, err := or.db.GetConnection().Exec(ctx,
		"insert into orders(number, status, accrual, uploaded_at, username) values($1, $2, $3, $4, $5) on conflict (number) do update set status = excluded.status, uploaded_at = excluded.uploaded_at",
		order.Number, order.Status, order.Accrual, order.UploadedAt, username,
	)
	if err != nil {
		log.Printf("failed to create order: %v", err)
		return errors.Errorf("failed to insert order '%v': %v", order, err)
	}
	return nil
}

func (or orderRepository) UpdateOrderStatus(ctx context.Context, number int, newStatus model.OrderStatus) error {
	log.Printf("Persisting update Order %d to %s", number, newStatus)
	_, err := or.db.GetConnection().Exec(ctx,
		"update orders set status = $1 where number = $2", newStatus, number,
	)
	if err != nil {
		log.Printf("failed to update order status '%d' to '%s': %v", number, newStatus, err)
		return errors.Errorf("failed to update order status'%d' to '%s': %v", number, newStatus, err)
	}
	return nil
}

func (or orderRepository) IsOrderExist(ctx context.Context, orderNumber int) (bool, error) {
	log.Printf("Checking %d Order existing", orderNumber)
	var orderID int
	row := or.db.GetConnection().QueryRow(ctx, "select number from orders where number = $1", fmt.Sprintf("%d", orderNumber))
	err := row.Scan(&orderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (or orderRepository) GetOrders(ctx context.Context, username string) ([]*model.Order, error) {
	log.Printf("Getting Orders for %s", username)
	var orders []*model.Order
	rows, err := or.db.GetConnection().Query(ctx, "select number, status, accrual, uploaded_at from orders where username = $1", username)
	if err != nil {
		log.Printf("failed to get orders for '%s': %v", username, err)
		return nil, errors.Errorf("failed to get orders for '%s': %v", username, err)
	}
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
