package usecase

import (
	"context"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountUsecase struct {
	repo AccountRepository
}

var (
	jwtKey = []byte("my_secret_key")
)

func NewAccountUsecase(repo AccountRepository) *AccountUsecase {
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

	expirationTime := time.Date(2023, 12, 31, 20, 34, 58, 651387237, time.UTC)
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
