package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WalletUsecase struct {
	repo WalletRepository
}

const (
	EnabledStatus  = 1
	DisabledStatus = 2
)

func NewWalletUsecase(repo WalletRepository) *WalletUsecase {
	return &WalletUsecase{
		repo: repo,
	}
}

func (h *WalletUsecase) Enabled(ctx context.Context, customerXid string) (*model.Wallet, error) {
	wallet, err := h.repo.GetWalletByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}

	if wallet != nil {
		if wallet.Status == EnabledStatus {
			return wallet, errors.New("Already enabled")
		}

		wallet.Status = EnabledStatus

		wallet, err = h.repo.UpdateStatusWallet(ctx, wallet)
		if err != nil {
			return nil, err
		}
	} else {
		payload := &model.Wallet{
			OwnedBy:   customerXid,
			Status:    EnabledStatus,
			EnabledAt: time.Now(),
		}

		wallet, err = h.repo.CreateWallet(ctx, payload)
		if err != nil {
			return nil, err
		}
	}

	wallet.StatusMessage = "enabled"

	return wallet, err
}

func (h *WalletUsecase) Disable(ctx context.Context, customerXid string) (*model.Wallet, error) {
	wallet, err := h.repo.GetWalletByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New("Wallet not enable")
	} else {
		if wallet.Status == DisabledStatus {
			return nil, errors.New("Wallet already disabled")
		}
	}

	wallet.Status = DisabledStatus

	wallet, err = h.repo.UpdateStatusWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	wallet.StatusMessage = "disabled"

	return wallet, err
}

func (h *WalletUsecase) GetWallet(ctx context.Context, customerXid string) (*model.Wallet, error) {
	infoWallet, err := h.repo.GetWalletByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}

	if infoWallet == nil {
		return nil, errors.New("Wallet not enable")
	}

	if infoWallet.Status == DisabledStatus {
		return nil, errors.New("Disabled")
	} else {
		infoWallet.StatusMessage = "enabled"
	}

	return infoWallet, err
}
