package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type DepositRepository interface {
	GetDepositByReferenceId(ctx context.Context, referenceId string) (*model.Deposit, error)
	CreateDeposit(ctx context.Context, deposit *model.Deposit, wallet *model.Wallet) (*model.Deposit, error)
}
