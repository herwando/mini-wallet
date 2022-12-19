package handler

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	mockUC "github.com/herwando/mini-wallet/module/wallet/handler/_mocks"
	"github.com/shopspring/decimal"
)

func TestHandlerDeposit_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDepositUC := mockUC.NewMockDepositUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockCustomerXid := "ea0212d3-abd6-406f-8c67-868e814a2436"
	mockToken := "Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl94aWQiOiJlYTAyMTJkMy1hYmQ2LTQwNmYtOGM2Ny04NjhlODE0YTI0MzYiLCJleHAiOjE2NzE0NjQwNDV9.hWYiJlo0bZiv3wxsr_6GOnIonmYjfex5gZR7ecbwB3U"
	mockCtx = context.WithValue(mockCtx, "AuthDetail", mockCustomerXid)
	mockCtxEmpty := context.TODO()
	now := time.Now()
	mockDeposit := &model.Deposit{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		DepositedBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		DepositedAt: now,
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
					var param = url.Values{}
					param.Add("reference_id", "ea0212d3-abd6-406f-8c67-868e814a2436")
					param.Add("amount", "1000")
					req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(param.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockDepositUC.EXPECT().CreateDeposit(gomock.Any(), mockCustomerXid, gomock.Any()).Return(mockDeposit, nil)
				writerWriteDataAccepted = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed params empty reference",
			args: args{
				r: func() *http.Request {
					var param = url.Values{}
					param.Add("amount", "1000")
					req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(param.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
					var param = url.Values{}
					param.Add("reference_id", "ea0212d3-abd6-406f-8c67-868e814a2436")
					param.Add("amount", "0")
					req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(param.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
					var param = url.Values{}
					param.Add("reference_id", "ea0212d3-abd6-406f-8c67-868e814a2436")
					param.Add("amount", "1000")
					req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(param.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
					var param = url.Values{}
					param.Add("reference_id", "ea0212d3-abd6-406f-8c67-868e814a2436")
					param.Add("amount", "1000")
					req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader(param.Encode()))
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockDepositUC.EXPECT().CreateDeposit(gomock.Any(), mockCustomerXid, gomock.Any()).Return(nil, mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewDepositHandler(mockDepositUC)
			h.Create(tt.args.w, tt.args.r)
		})
	}
}
