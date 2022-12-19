package handler

import (
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
)

type WalletHandler struct {
	usecase WalletUsecase
}

func NewWalletHandler(usecase WalletUsecase) *WalletHandler {
	return &WalletHandler{
		usecase: usecase,
	}
}

func (h *WalletHandler) Enabled(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXid, err := middlewares.GetAuthDetailFromContext(ctx)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Token", err.Error()))
		return
	}

	wallet, err := h.usecase.Enabled(ctx, customerXid)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Wallet", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   wallet,
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}

func (h *WalletHandler) Disable(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	ctx := r.Context()
	var payload model.PayloadDisable
	isDisable := r.Form.Get("is_disabled")
	if isDisable == "true" {
		payload.IsDisabled = true
	}

	if !payload.IsDisabled {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Params is_disabled empty"))
		return
	}

	customerXid, err := middlewares.GetAuthDetailFromContext(ctx)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Token", err.Error()))
		return
	}

	wallet, err := h.usecase.Disable(ctx, customerXid)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Wallet", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   wallet,
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerXid, err := middlewares.GetAuthDetailFromContext(ctx)
	if err != nil {
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Token", err.Error()))
		return
	}

	wallet, err := h.usecase.GetWallet(ctx, customerXid)
	if err != nil {
		if err.Error() == "Disabled" {
			writerWriteJSONAPIError(ctx, w, commonerr.SetNewNotFound(err.Error()))
			return
		}
		writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnprocessableEntity("Wallet", err.Error()))
		return
	}

	response := BasicResponse{
		Data:   wallet,
		Status: "success",
	}

	writerWriteData(ctx, w, response)
}
