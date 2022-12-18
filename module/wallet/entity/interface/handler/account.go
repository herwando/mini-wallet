package handler

import (
	"net/http"
)

type AccountHandler interface {
	Init(w http.ResponseWriter, r *http.Request)
}
