package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WalletUsecase interface {
	Enabled(ctx context.Context, customerXid string) (*model.Wallet, error)
	Disable(ctx context.Context, customerXid string) (*model.Wallet, error)
	GetWallet(ctx context.Context, customerXid string) (*model.Wallet, error)
}
