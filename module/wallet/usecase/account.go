package usecase

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/repository"
)

type AccountUsecase struct {
	repo *repository.AccountRepository
}

var (
	jwtKey = []byte("my_secret_key")
)

func NewAccountUsecase(repo *repository.AccountRepository) *AccountUsecase {
	return &AccountUsecase{
		repo: repo,
	}
}

func (h *AccountUsecase) Init(ctx context.Context, payload model.Account) (string, error) {
	exist, err := h.repo.ExistAccountByCustomerXid(ctx, payload.CustomerXid)
	if err != nil {
		return "", err
	}

	if !exist {
		err := h.repo.CreateAccount(ctx, &payload)
		if err != nil {
			return "", err
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &model.Claims{
		CustomerXid: payload.CustomerXid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}
