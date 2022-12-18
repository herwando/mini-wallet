package handler

import (
	"encoding/json"
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
)

type DepositHandler struct {
	usecase *usecase.DepositUsecase
}

func NewDepositHandler(usecase *usecase.DepositUsecase) *DepositHandler {
	return &DepositHandler{
		usecase: usecase,
	}
}

func (h *DepositHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.PayloadDeposit
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Body not completed"))
		return
	}

	if payload.ReferenceID == "" {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Params reference_id empty"))
		return
	}

	if !payload.Amount.IsPositive() {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Params amount not valid"))
		return
	}

	customerXid, err := middlewares.GetAuthDetailFromContext(ctx)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Token", err.Error()))
		return
	}

	deposit, err := h.usecase.CreateDeposit(ctx, customerXid, payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Deposit", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   deposit,
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}
