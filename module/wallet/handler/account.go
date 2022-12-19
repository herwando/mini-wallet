package handler

import (
	"encoding/json"
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
	ctx := r.Context()
	var payload model.Account
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Body not completed"))
		return
	}

	if payload.CustomerXid == "" {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Params customer_xid empty"))
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
