package repository

import "context"

type OrderRepository interface {
	CreateOrder(ctx context.Context, username string, orderId int) error
	GetOrders(ctx context.Context, username string) ([]int, error)
}

type orderRepository struct {
	db DbProvider
}

func NewOrderRepository(db DbProvider) OrderRepository {
	return &orderRepository{db}
}

func (or orderRepository) CreateOrder(ctx context.Context, username string, orderId int) error {
	_, err := or.db.GetConnection().Exec(ctx, "insert into orders(username, order_id) values($1, $2)", username, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (or orderRepository) GetOrders(ctx context.Context, username string) ([]int, error) {
	var orderIds []int
	rows, err := or.db.GetConnection().Query(ctx, "select order_id from orders where username = $1", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var orderId int
		err := rows.Scan(&orderId)
		if err != nil {
			return nil, err
		}
		orderIds = append(orderIds, orderId)
	}
	return orderIds, nil
}
