package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountRepository interface {
	ExistAccountByCustomerXid(ctx context.Context, customerXid string) (bool, error)
	CreateAccount(ctx context.Context, account *model.Account) error
}

type DepositRepository interface {
	GetDepositByReferenceId(ctx context.Context, referenceId string) (*model.Deposit, error)
	CreateDeposit(ctx context.Context, deposit *model.Deposit, wallet *model.Wallet) (*model.Deposit, error)
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	UpdateStatusWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	GetWalletByCustomerXid(ctx context.Context, customerXid string) (*model.Wallet, error)
}

type WithdrawalRepository interface {
	GetWithdrawalByReferenceId(ctx context.Context, referenceId string) (*model.Withdrawal, error)
	CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal, wallet *model.Wallet) (*model.Withdrawal, error)
}
