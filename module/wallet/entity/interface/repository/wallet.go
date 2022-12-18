package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	UpdateStatusWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)
	GetWalletByCustomerXid(ctx context.Context, customerXid string) (*model.Wallet, error)
}
