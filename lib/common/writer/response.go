package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/pkg/errors"
)

type errorCtx struct {
	err error
}

var (
	ok          = []byte(`{"data":"ok"}`)
	ErrorCtxKey = errorCtx{}
)

// WriteOK is a helper around Response.Write with status OK
func WriteOK(ctx context.Context, w http.ResponseWriter, data interface{}) {
	write(ctx, w, data, http.StatusOK)
}

func WriteAccepted(ctx context.Context, w http.ResponseWriter, data interface{}) {
	write(ctx, w, data, http.StatusAccepted)
}

// WriteStrOK is a helper around Response.Write with status OK
func WriteStrOK(ctx context.Context, w http.ResponseWriter) {
	set(ctx, w, ok, http.StatusOK)
}

// WriteJSONAPIError is a helper
func WriteJSONAPIError(ctx context.Context, w http.ResponseWriter, err error) {
	switch errCause := errors.Cause(err).(type) {
	case *commonerr.ErrorMessage:
		write(ctx, w, errCause, errCause.Code)
	default:
		write(ctx, w, commonerr.ErrorMessage{
			ErrorList: commonerr.SetNewInternalError().GetListError(),
		}, http.StatusInternalServerError)
	}
}

func write(ctx context.Context, w http.ResponseWriter, data interface{}, status int) {
	datab, err := json.Marshal(data)
	if err != nil {
		datab = []byte(`{"error_list":[{"error_name": "internal", "error_description": "Internal Server Error"}]}`)
	}
	set(ctx, w, datab, status)
}

func set(ctx context.Context, w http.ResponseWriter, datab []byte, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(datab)
	if err != nil {
		fmt.Println("[HTTP]", err)
	}
}
