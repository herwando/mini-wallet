package handler

import (
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/writer"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
)

var (
	writerWriteStrOK        = writer.WriteStrOK
	writerWriteData         = writer.WriteOK
	writerWriteJSONAPIError = writer.WriteJSONAPIError
)

type WalletHandler struct {
	usecase *usecase.WalletUsecase
}

func NewWalletHandler(usecase *usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		usecase: usecase,
	}
}

func (h *WalletHandler) Ping(w http.ResponseWriter, r *http.Request) {
	writerWriteStrOK(r.Context(), w)
}
