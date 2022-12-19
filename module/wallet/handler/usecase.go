package handler

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountUsecase interface {
	Init(ctx context.Context, payload model.Account) (string, error)
}

type DepositUsecase interface {
	CreateDeposit(ctx context.Context, customerXid string, payload model.PayloadDeposit) (*model.Deposit, error)
}

type WalletUsecase interface {
	Enabled(ctx context.Context, customerXid string) (*model.Wallet, error)
	Disable(ctx context.Context, customerXid string) (*model.Wallet, error)
	GetWallet(ctx context.Context, customerXid string) (*model.Wallet, error)
}

type WithdrawalUsecase interface {
	CreateWithdrawal(ctx context.Context, customerXid string, payload model.PayloadWithdrawal) (*model.Withdrawal, error)
}
