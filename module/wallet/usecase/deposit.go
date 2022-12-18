package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/repository"
)

type DepositUsecase struct {
	repo       *repository.DepositRepository
	walletRepo *repository.WalletRepository
}

const (
	SuccessStatus = 1
)

func NewDepositUsecase(repo *repository.DepositRepository, walletRepo *repository.WalletRepository) *DepositUsecase {
	return &DepositUsecase{
		repo:       repo,
		walletRepo: walletRepo,
	}
}

func (h *DepositUsecase) CreateDeposit(ctx context.Context, customerXid string, payload model.PayloadDeposit) (*model.Deposit, error) {
	wallet, err := h.walletRepo.GetWalletByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New("Wallet not enable")
	} else {
		if wallet.Status == DisabledStatus {
			return nil, errors.New("Wallet disabled")
		}
	}
	wallet.Balance = wallet.Balance.Add(payload.Amount)

	deposit, err := h.repo.GetDepositByReferenceId(ctx, payload.ReferenceID)
	if err != nil {
		return nil, err
	}

	if deposit != nil {
		return nil, errors.New("Reference id already used")
	}

	data := &model.Deposit{
		DepositedBy: customerXid,
		Status:      SuccessStatus,
		DepositedAt: time.Now(),
		Amount:      payload.Amount,
		ReferenceID: payload.ReferenceID,
	}

	deposit, err = h.repo.CreateDeposit(ctx, data, wallet)
	if err != nil {
		return nil, err
	}
	deposit.StatusMessage = "success"

	return deposit, err
}
