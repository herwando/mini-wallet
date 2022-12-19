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

func TestHandlerWallet_Enabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletUC := mockUC.NewMockWalletUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockCustomerXid := "ea0212d3-abd6-406f-8c67-868e814a2436"
	mockCtx = context.WithValue(mockCtx, "AuthDetail", mockCustomerXid)
	mockToken := "Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl94aWQiOiJlYTAyMTJkMy1hYmQ2LTQwNmYtOGM2Ny04NjhlODE0YTI0MzYiLCJleHAiOjE2NzE0NjQwNDV9.hWYiJlo0bZiv3wxsr_6GOnIonmYjfex5gZR7ecbwB3U"
	now := time.Now()
	mockWalletEnable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCtxEmpty := context.TODO()

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
					req, _ := http.NewRequest(http.MethodPost, "", http.NoBody)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().Enabled(gomock.Any(), mockCustomerXid).Return(mockWalletEnable, nil)
				writerWriteData = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed token",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "", http.NoBody)
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
					req, _ := http.NewRequest(http.MethodPost, "", http.NoBody)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().Enabled(gomock.Any(), mockCustomerXid).Return(nil, mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewWalletHandler(mockWalletUC)
			h.Enabled(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlerWallet_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletUC := mockUC.NewMockWalletUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockCustomerXid := "ea0212d3-abd6-406f-8c67-868e814a2436"
	mockCtx = context.WithValue(mockCtx, "AuthDetail", mockCustomerXid)
	mockToken := "Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl94aWQiOiJlYTAyMTJkMy1hYmQ2LTQwNmYtOGM2Ny04NjhlODE0YTI0MzYiLCJleHAiOjE2NzE0NjQwNDV9.hWYiJlo0bZiv3wxsr_6GOnIonmYjfex5gZR7ecbwB3U"
	now := time.Now()
	mockWalletDisable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCtxEmpty := context.TODO()
	paramsReq := model.PayloadDisable{
		IsDisabled: true,
	}
	paramsEmpty := model.Wallet{}
	reqByte, _ := json.Marshal(&paramsReq)
	reqByte2, _ := json.Marshal(&paramsEmpty)

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
					req, _ := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(reqByte))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().Disable(gomock.Any(), mockCustomerXid).Return(mockWalletDisable, nil)
				writerWriteData = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed params diff model",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(reqByte2))
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
					req, _ := http.NewRequest(http.MethodPatch, "", http.NoBody)
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
					req, _ := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(reqByte))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().Disable(gomock.Any(), mockCustomerXid).Return(nil, mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewWalletHandler(mockWalletUC)
			h.Disable(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlerWallet_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletUC := mockUC.NewMockWalletUsecase(ctrl)
	defer ctrl.Finish()

	mockCtx := context.TODO()
	mockError := errors.New("fake error")
	mockCustomerXid := "ea0212d3-abd6-406f-8c67-868e814a2436"
	mockCtx = context.WithValue(mockCtx, "AuthDetail", mockCustomerXid)
	mockToken := "Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcl94aWQiOiJlYTAyMTJkMy1hYmQ2LTQwNmYtOGM2Ny04NjhlODE0YTI0MzYiLCJleHAiOjE2NzE0NjQwNDV9.hWYiJlo0bZiv3wxsr_6GOnIonmYjfex5gZR7ecbwB3U"
	now := time.Now()
	mockWalletEnable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCtxEmpty := context.TODO()

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
					req, _ := http.NewRequest(http.MethodGet, "", http.NoBody)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().GetWallet(gomock.Any(), mockCustomerXid).Return(mockWalletEnable, nil)
				writerWriteData = func(ctx context.Context, w http.ResponseWriter, data interface{}) {
				}
			},
		},
		{
			name: "Failed token",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "", http.NoBody)
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
					req, _ := http.NewRequest(http.MethodGet, "", http.NoBody)
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", mockToken)
					req = req.WithContext(mockCtx)
					return req
				}(),
			},
			patch: func() {
				mockWalletUC.EXPECT().GetWallet(gomock.Any(), mockCustomerXid).Return(nil, mockError)
				writerWriteJSONAPIError = func(ctx context.Context, w http.ResponseWriter, err error) {
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			h := NewWalletHandler(mockWalletUC)
			h.GetWallet(tt.args.w, tt.args.r)
		})
	}
}
