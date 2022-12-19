package handler

import (
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
	"github.com/shopspring/decimal"
)

type DepositHandler struct {
	usecase DepositUsecase
}

func NewDepositHandler(usecase DepositUsecase) *DepositHandler {
	return &DepositHandler{
		usecase: usecase,
	}
}

func (h *DepositHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	ctx := r.Context()
	var payload model.PayloadDeposit
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

	deposit, err := h.usecase.CreateDeposit(ctx, customerXid, payload)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Deposit", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   deposit,
		Status: "success",
	}

	writerWriteDataAccepted(ctx, w, response)
}
