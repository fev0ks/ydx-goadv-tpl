package repository

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, order *model.Order) error
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderID int) (*model.Order, error)
	GetOrders(ctx context.Context, userID int) ([]*model.Order, error)
}

type orderRepository struct {
	db DBProvider
}

func NewOrderRepository(db DBProvider) OrderRepository {
	return &orderRepository{db}
}

func (or orderRepository) InsertOrder(ctx context.Context, order *model.Order) error {
	log.Printf("Persisting Order: %v", order)
	conn, err := or.db.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Printf("failed to open order tx '%d': %v", order.UserID, err)
		return errors.Errorf("failed to order open tx '%d': %v", order.UserID, err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		"insert into orders(order_id, status, accrual, uploaded_at) values($1, $2, $3, $4)",
		order.Number, order.Status, order.Accrual, order.UploadedAt,
	)
	if err != nil {
		log.Printf("failed to create order: %v", err)
		return errors.Errorf("failed to insert order '%v': %v", order, err)
	}
	_, err = tx.Exec(ctx, "insert into user_orders(user_id, order_id) values($1, $2)", order.UserID, order.Number)
	if err != nil {
		log.Printf("failed to insert user '%d' - order '%d' relation: %v", order.UserID, order.Number, err)
		return errors.Errorf("failed to insert user '%d' - order '%d' relation: %v", order.UserID, order.Number, err)
	}
	err = or.updateUserBalance(ctx, tx, order)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (or orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	log.Printf("Persisting update Order %v", order)
	conn, err := or.db.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Printf("failed to open order tx '%d': %v", order.UserID, err)
		return errors.Errorf("failed to order open tx '%d': %v", order.UserID, err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "update orders set status = $1, accrual = $2 where order_id = $3", order.Status, order.Accrual, order.Number)
	if err != nil {
		log.Printf("failed to update order %v: %v", order, err)
		return errors.Errorf("failed to update order %v: %v", order, err)
	}
	err = or.updateUserBalance(ctx, tx, order)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (or orderRepository) updateUserBalance(ctx context.Context, tx pgx.Tx, order *model.Order) error {
	if order.Accrual != 0 && order.Status == model.ProcessedStatus {
		_, err := tx.Exec(ctx, "update user_balance set current = current + $1 where user_id = $2", order.Accrual, order.UserID)
		if err != nil {
			log.Printf("failed to update user balance %d: %v", order.UserID, err)
			return errors.Errorf("failed to update user balance %d: %v", order.UserID, err)
		}
	}
	return nil
}

func (or orderRepository) GetOrder(ctx context.Context, orderID int) (*model.Order, error) {
	log.Printf("Checking '%d' Order existing", orderID)
	order := &model.Order{}
	conn, err := or.db.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx,
		"select uo.user_id, o.order_id, o.status, o.accrual, o.uploaded_at from orders o join user_orders uo on uo.order_id = o.order_id and o.order_id = $1", orderID)
	err = row.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("'%d' Order doesn't exist", orderID)
			return nil, nil
		}
		return nil, err
	}
	log.Printf("'%d' Order exists", orderID)
	return order, nil
}

func (or orderRepository) GetOrders(ctx context.Context, userID int) ([]*model.Order, error) {
	log.Printf("Getting Orders for userID: %d", userID)
	var orders []*model.Order
	conn, err := or.db.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		"select uo.user_id, o.order_id, o.status, o.accrual, o.uploaded_at from orders o join user_orders uo on uo.user_id = $1 and o.order_id = uo.order_id", userID)
	if err != nil {
		log.Printf("failed to get orders for '%d': %v", userID, err)
		return nil, errors.Errorf("failed to get orders for '%d': %v", userID, err)
	}
	defer rows.Close()
	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			log.Printf("failed to read orders for '%d': %v", userID, err)
			return nil, errors.Errorf("failed to read orders for '%d': %v", userID, err)
		}
		orders = append(orders, order)
	}
	log.Printf("Got %d Orders for userID: %d", len(orders), userID)
	return orders, nil
}
