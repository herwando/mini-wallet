package handler

import (
	"encoding/json"
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
)

type WithdrawalHandler struct {
	usecase WithdrawalUsecase
}

func NewWithdrawalHandler(usecase WithdrawalUsecase) *WithdrawalHandler {
	return &WithdrawalHandler{
		usecase: usecase,
	}
}

func (h *WithdrawalHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.PayloadWithdrawal
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

	withdrawal, err := h.usecase.CreateWithdrawal(ctx, customerXid, payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Withdrawal", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   withdrawal,
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}
