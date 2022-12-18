package usecase

import (
	"context"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountUsecase interface {
	Init(ctx context.Context, payload model.Account) (string, error)
}
