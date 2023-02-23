package repository

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (*model.Balance, error)
	BalanceWithdraw(ctx context.Context, userID int, withdraw *model.Withdraw) error
	GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdraw, error)
}

type balanceRepository struct {
	db DBProvider
}

func NewBalancewRepository(db DBProvider) BalanceRepository {
	return &balanceRepository{db}
}

func (br *balanceRepository) GetBalance(ctx context.Context, userID int) (*model.Balance, error) {
	conn, err := br.db.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	result := conn.QueryRow(ctx,
		"select current, withdraw from user_balance where user_id = $1",
		userID)
	balance := &model.Balance{}
	err = result.Scan(&balance.Current, &balance.Withdraw)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get balance of '%d': %v", userID, err)
	}
	return balance, nil
}

func (br *balanceRepository) BalanceWithdraw(ctx context.Context, userID int, withdraw *model.Withdraw) error {
	conn, err := br.db.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Printf("failed to open balance tx '%d': %v", userID, err)
		return errors.Errorf("failed to open balance tx '%d': %v", userID, err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "update user_balance set current = current - $1, withdraw = withdraw + $2 where user_id = $3", withdraw.Sum, withdraw.Sum, userID)
	if err != nil {
		log.Printf("failed to withdraw for '%d': %v", userID, err)
		return errors.Errorf("failed to withdraw for '%d': %v", userID, err)
	}
	withdrawRow := tx.QueryRow(ctx,
		"insert into withdraws(order_id, sum, processed_at) values ($1, $2, $3) returning withdraw_id",
		withdraw.Order, withdraw.Sum, withdraw.ProcessedAt)
	var WithdrawID int
	err = withdrawRow.Scan(&WithdrawID)
	if err != nil {
		log.Printf("failed to get withdrawID '%d': %v", userID, err)
		return errors.Errorf("failed to get withdrawID '%d': %v", userID, err)
	}
	_, err = tx.Exec(ctx,
		"insert into user_withdraws(user_id, withdraw_id) values ($1, $2)",
		userID, WithdrawID)
	if err != nil {
		log.Printf("failed to insert '%d' user - '%d' withdraw: %v", userID, WithdrawID, err)
		return errors.Errorf("failed to insert '%d' user - '%d' withdraw: %v", userID, WithdrawID, err)
	}
	tx.Commit(ctx)
	return nil
}

func (br *balanceRepository) GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdraw, error) {
	log.Printf("Getting withdraws for userID: %d", userID)
	var withdraws []*model.Withdraw
	conn, err := br.db.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		"select w.order_id, w.sum, w.processed_at from withdraws w join user_withdraws u on u.withdraw_id = w.withdraw_id and u.user_id = $1", userID)
	if err != nil {
		log.Printf("failed to get withdraws for '%d': %v", userID, err)
		return nil, errors.Errorf("failed to get withdraws for '%d': %v", userID, err)
	}
	defer rows.Close()
	for rows.Next() {
		withdraw := &model.Withdraw{}
		err := rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			log.Printf("failed to get withdraws for '%d': %v", userID, err)
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}
