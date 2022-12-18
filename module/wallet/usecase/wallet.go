package usecase

import (
	"github.com/herwando/mini-wallet/module/wallet/repository"
)

type WalletUsecase struct {
	repo *repository.WalletRepository
}

func NewWalletUsecase(repo *repository.WalletRepository) *WalletUsecase {
	return &WalletUsecase{
		repo: repo,
	}
}
