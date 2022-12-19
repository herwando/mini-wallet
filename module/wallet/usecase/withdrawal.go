package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WithdrawalUsecase struct {
	repo       WithdrawalRepository
	walletRepo WalletRepository
}

func NewWithdrawalUsecase(repo WithdrawalRepository, walletRepo WalletRepository) *WithdrawalUsecase {
	return &WithdrawalUsecase{
		repo:       repo,
		walletRepo: walletRepo,
	}
}

func (h *WithdrawalUsecase) CreateWithdrawal(ctx context.Context, customerXid string, payload model.PayloadWithdrawal) (*model.Withdrawal, error) {
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

		if wallet.Balance.LessThan(payload.Amount) {
			return nil, errors.New("Wallet balance not enough")
		}
	}
	wallet.Balance = wallet.Balance.Sub(payload.Amount)

	withdrawal, err := h.repo.GetWithdrawalByReferenceId(ctx, payload.ReferenceID)
	if err != nil {
		return nil, err
	}

	if withdrawal != nil {
		return nil, errors.New("Reference id already used")
	}

	data := &model.Withdrawal{
		WithdrawnBy: customerXid,
		Status:      SuccessStatus,
		WithdrawnAt: time.Now(),
		Amount:      payload.Amount,
		ReferenceID: payload.ReferenceID,
	}

	withdrawal, err = h.repo.CreateWithdrawal(ctx, data, wallet)
	if err != nil {
		return nil, err
	}
	withdrawal.StatusMessage = "success"

	return withdrawal, err
}
