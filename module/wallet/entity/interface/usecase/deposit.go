package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type DepositUsecase interface {
	CreateDeposit(ctx context.Context, customerXid string, payload model.PayloadDeposit) (*model.Deposit, error)
}
