package handler

import (
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
	"github.com/shopspring/decimal"
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
	r.ParseMultipartForm(32 << 20)
	ctx := r.Context()
	var payload model.PayloadWithdrawal
	payload.ReferenceID = r.Form.Get("reference_id")
	amount, err := decimal.NewFromString(r.Form.Get("amount"))
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Params amount not valid"))
		return
	}

	payload.Amount = amount

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

	writerWriteDataAccepted(ctx, w, response)
}
