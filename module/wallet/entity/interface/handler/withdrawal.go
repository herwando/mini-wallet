package handler

import (
	"net/http"
)

type WithdrawalHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}
