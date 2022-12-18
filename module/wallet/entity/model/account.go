package model

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type Account struct {
	CustomerXid string `json:"customer_xid"`
}

type Claims struct {
	CustomerXid string `json:"customer_xid"`
	jwt.RegisteredClaims
}

type AccountResponse struct {
	Token string `json:"token"`
}
