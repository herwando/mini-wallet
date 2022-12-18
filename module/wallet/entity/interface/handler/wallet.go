package handler

import (
	"net/http"
)

type WalletHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}
