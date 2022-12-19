package handler

import (
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountHandler struct {
	usecase AccountUsecase
}

func NewAccountHandler(usecase AccountUsecase) *AccountHandler {
	return &AccountHandler{
		usecase: usecase,
	}
}

func (h *AccountHandler) Init(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	ctx := r.Context()
	var payload model.Account
	payload.CustomerXid = r.Form.Get("customer_xid")

	if payload.CustomerXid == "" {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("customer_xid", "Missing data for required field"))
		return
	}

	tokenString, err := h.usecase.Init(ctx, payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Account", err.Error()))
		return
	}

	response := BasicResponse{
		Data: model.AccountResponse{
			Token: tokenString,
		},
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}
