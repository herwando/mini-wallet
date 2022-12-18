package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WithdrawalRepository interface {
	GetWithdrawalByReferenceId(ctx context.Context, referenceId string) (*model.Withdrawal, error)
	CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal, wallet *model.Wallet) (*model.Withdrawal, error)
}
