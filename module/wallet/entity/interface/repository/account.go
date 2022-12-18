package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountRepository interface {
	ExistAccountByCustomerXid(ctx context.Context, customerXid string) (bool, error)
	CreateAccount(ctx context.Context, account *model.Account) error
}
