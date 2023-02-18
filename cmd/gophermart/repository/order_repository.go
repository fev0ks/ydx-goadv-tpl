package repository

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, userID int, order *model.Order) error
	UpdateOrder(ctx context.Context, userID int, order *model.Order) error
	IsOrderExist(ctx context.Context, orderID int) (bool, error)
	GetOrders(ctx context.Context, userID int) ([]*model.Order, error)
}

type orderRepository struct {
	db DBProvider
}

func NewOrderRepository(db DBProvider) OrderRepository {
	return &orderRepository{db}
}

func (or orderRepository) InsertOrder(ctx context.Context, userID int, order *model.Order) error {
	log.Printf("Persisting Order: %v", order)
	tx, err := or.db.GetConnection().Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		log.Printf("failed to open Tx: %v", err)
		return err
	}
	_, err = tx.Exec(ctx,
		"insert into orders(order_id, status, accrual, uploaded_at) values($1, $2, $3, $4)",
		order.Number, order.Status, order.Accrual, order.UploadedAt,
	)
	if err != nil {
		log.Printf("failed to create order: %v", err)
		return errors.Errorf("failed to insert order '%v': %v", order, err)
	}
	_, err = tx.Exec(ctx, "insert into user_orders(user_id, order_id) values($1, $2)", userID, order.Number)
	if err != nil {
		log.Printf("failed to insert user %d - order %d relation: %v", userID, order.Number, err)
		return errors.Errorf("failed to insert user %d - order %d relation: %v", userID, order.Number, err)
	}
	err = or.updateUserBalance(ctx, tx, userID, order)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (or orderRepository) UpdateOrder(ctx context.Context, userID int, order *model.Order) error {
	log.Printf("Persisting update Order %v", order)
	tx, err := or.db.GetConnection().Begin(ctx)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "update orders set status = $1, accrual = $2 where order_id = $2", order.Status, order.Accrual, order.Number)
	if err != nil {
		log.Printf("failed to update order %v: %v", order, err)
		return errors.Errorf("failed to update order %v: %v", order, err)
	}
	err = or.updateUserBalance(ctx, tx, userID, order)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (or orderRepository) updateUserBalance(ctx context.Context, tx pgx.Tx, userID int, order *model.Order) error {
	if order.Accrual != 0 {
		_, err := tx.Exec(ctx, "update user_balance set current = current + $1 where user_id = $2", order.Accrual, userID)
		if err != nil {
			log.Printf("failed to update user balance %d: %v", userID, err)
			return errors.Errorf("failed to update user balance %d: %v", userID, err)
		}
	}
	return nil
}

func (or orderRepository) IsOrderExist(ctx context.Context, orderID int) (bool, error) {
	log.Printf("Checking %d Order existing", orderID)
	row := or.db.GetConnection().QueryRow(ctx, "select order_id from orders where order_id = $1", orderID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (or orderRepository) GetOrders(ctx context.Context, userID int) ([]*model.Order, error) {
	log.Printf("Getting Orders for userID: %d", userID)
	var orders []*model.Order
	rows, err := or.db.GetConnection().Query(ctx,
		"select o.order_id, o.status, o.accrual, o.uploaded_at from orders o join user_orders u on u.user_id = $1 and o.order_id = u.order_id", userID)
	if err != nil {
		log.Printf("failed to get orders for '%d': %v", userID, err)
		return nil, errors.Errorf("failed to get orders for '%d': %v", userID, err)
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
