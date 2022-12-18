package handler

import (
	"net/http"
)

type DepositHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}
