package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/internal/model"
	"github.com/fev0ks/ydx-goadv-tpl/internal/repository"
	"time"
)

type BalanceService interface {
	GetBalance(ctx context.Context, userID int) (*model.Balance, error)
	BalanceWithdraw(ctx context.Context, userID int, withdrawRequest *model.WithdrawRequest) error
	GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdraw, error)
}

type balanceService struct {
	balanceRepo repository.BalanceRepository
}

func NewBalanceService(balanceRepo repository.BalanceRepository) BalanceService {
	return &balanceService{balanceRepo}
}

func (bs *balanceService) GetBalance(ctx context.Context, userID int) (*model.Balance, error) {
	return bs.balanceRepo.GetBalance(ctx, userID)
}

func (bs *balanceService) BalanceWithdraw(ctx context.Context, userID int, withdrawRequest *model.WithdrawRequest) error {
	return bs.balanceRepo.BalanceWithdraw(ctx, userID, &model.Withdraw{
		Order:       withdrawRequest.Order,
		Sum:         withdrawRequest.Sum,
		ProcessedAt: time.Now(),
	})
}

func (bs *balanceService) GetWithdrawals(ctx context.Context, userID int) ([]*model.Withdraw, error) {
	return bs.balanceRepo.GetWithdrawals(ctx, userID)
}
