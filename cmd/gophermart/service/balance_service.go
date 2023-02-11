package service

type BalanceService interface {
	GetBalance()
	BalanceWithdraw()
	GetWithdrawsHistory()
}
