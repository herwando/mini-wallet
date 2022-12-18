package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WithdrawalUsecase interface {
	CreateWithdrawal(ctx context.Context, customerXid string, payload model.PayloadWithdrawal) (*model.Withdrawal, error)
}
