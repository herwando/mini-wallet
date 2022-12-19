package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	mockUC "github.com/herwando/mini-wallet/module/wallet/handler/_mocks"
	"github.com/shopspring/decimal"
)

func TestHandlerWithdrawal_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWithdrawalUC := mockUC.NewMockWithdrawalUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockCustomerXid := "ea0212d3-abd6-406f-8c67-868e814a2436"
	mockToken := "Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl94aWQiOiJlYTAyMTJkMy1hYmQ2LTQwNmYtOGM2Ny04NjhlODE0YTI0MzYiLCJleHAiOjE2NzE0NjQwNDV9.hWYiJlo0bZiv3wxsr_6GOnIonmYjfex5gZR7ecbwB3U"
	paramsReq := model.PayloadWithdrawal{
		ReferenceID: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Amount:      decimal.NewFromInt(1000),
	}
	paramsEmpty := model.Wallet{}
	paramsEmptyReference := model.PayloadWithdrawal{
		ReferenceID: "",
	}
	paramsAmountZero := model.PayloadWithdrawal{
		ReferenceID: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Amount:      decimal.NewFromInt(0),
	}
	mockCtx = context.WithValue(mockCtx, "AuthDetail", mockCustomerXid)
	mockCtxEmpty := context.TODO()
	reqByte, _ := json.Marshal(&paramsReq)
	reqByte2, _ := json.Marshal(&paramsEmpty)
	reqByte3, _ := json.Marshal(&paramsEmptyReference)
	reqByte4, _ := json.Marshal(&paramsAmountZero)
	now := time.Now()
	mockWithdrawal := &model.Withdrawal{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		WithdrawnBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		WithdrawnAt: now,
		Amount:      decimal.NewFromInt(1000),
		ReferenceID: "81401b03-60e0-4f20-afc6-419b3773e7b3",
	}

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
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWithdrawalUC.EXPECT().CreateWithdrawal(gomock.Any(), mockCustomerXid, gomock.Any()).Return(mockWithdrawal, nil)
				writerWriteData = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed params diff model",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte2))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
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
			name: "Failed params empty reference",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte3))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
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
			name: "Failed params zero amount",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte4))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
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
			name: "Failed token",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqByte))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtxEmpty)
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
					req.Header.Set("Authorization", mockToken)
					return req
				}(),
			},
			patch: func() {
				mockWithdrawalUC.EXPECT().CreateWithdrawal(gomock.Any(), mockCustomerXid, gomock.Any()).Return(nil, mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewWithdrawalHandler(mockWithdrawalUC)
			h.Create(tt.args.w, tt.args.r)
		})
	}
}
