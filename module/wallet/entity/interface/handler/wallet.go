package handler

import (
	"net/http"
)

type WalletHandler interface {
	Enabled(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
}
