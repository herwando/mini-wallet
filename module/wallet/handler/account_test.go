package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	mockUC "github.com/herwando/mini-wallet/module/wallet/handler/_mocks"
)

func TestHandlerAccount_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountUC := mockUC.NewMockAccountUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockToken := "exampletoken"
	paramsReq := model.Account{
		CustomerXid: "ea0212d3-abd6-406f-8c67-868e814a2436",
	}
	paramsEmpty := model.Wallet{}
	paramsEmptyString := model.Account{
		CustomerXid: "",
	}
	reqByte, _ := json.Marshal(&paramsReq)
	reqByte2, _ := json.Marshal(&paramsEmpty)
	reqByte3, _ := json.Marshal(&paramsEmptyString)

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name  string
		args  args
		patch func()
	}{
		{
			name: "Success",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte))
					req.Header.Set("Content-Type", "application/json")
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockAccountUC.EXPECT().Init(gomock.Any(), gomock.Any()).Return(mockToken, nil)
				writerWriteData = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed params empty",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte2))
					req.Header.Set("Content-Type", "application/json")
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
		{
			name: "Failed params empty string",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte3))
					req.Header.Set("Content-Type", "application/json")
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
		{
			name: "Failed usecase",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte))
					req.Header.Set("Content-Type", "application/json")
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockAccountUC.EXPECT().Init(gomock.Any(), gomock.Any()).Return("", mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewAccountHandler(mockAccountUC)
			h.Init(tt.args.w, tt.args.r)
		})
	}
}
