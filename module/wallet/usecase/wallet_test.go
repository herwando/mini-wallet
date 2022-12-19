package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
	mockDB "github.com/herwando/mini-wallet/module/wallet/usecase/_mocks"
	"github.com/shopspring/decimal"
)

func TestUsecaseWallet_GetWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletDB := mockDB.NewMockWalletRepository(ctrl)
	defer ctrl.Finish()
	now := time.Now()
	mockCtx := context.Background()
	mockError := errors.New("fake error")
	mockWallet := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockWalletDisable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCustomerXid := mockWallet.OwnedBy

	type args struct {
		ctx         context.Context
		customerXid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "Success",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: false,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
			},
		},
		{
			name: "Failed on GetWalletByCustomerXid",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, mockError)
			},
		},
		{
			name: "Failed wallet empty",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, nil)
			},
		},
		{
			name: "Failed wallet disabled",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWalletDisable, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewWalletUsecase(mockWalletDB)
			_, err := uc.GetWallet(tt.args.ctx, tt.args.customerXid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseWallet_Enabled(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletDB := mockDB.NewMockWalletRepository(ctrl)
	defer ctrl.Finish()
	now := time.Now()
	mockCtx := context.Background()
	mockError := errors.New("fake error")
	mockWallet := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockWalletDisable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCustomerXid := mockWallet.OwnedBy

	type args struct {
		ctx         context.Context
		customerXid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "Success with create",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: false,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, nil)
				mockWalletDB.EXPECT().CreateWallet(gomock.Any(), gomock.Any()).Return(mockWallet, nil)
			},
		},
		{
			name: "Failed on CreateWallet",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, nil)
				mockWalletDB.EXPECT().CreateWallet(gomock.Any(), gomock.Any()).Return(nil, mockError)
			},
		},
		{
			name: "Failed on UpdateStatusWallet",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWalletDisable, nil)
				mockWalletDB.EXPECT().UpdateStatusWallet(gomock.Any(), gomock.Any()).Return(nil, mockError)
			},
		},
		{
			name: "Failed wallet already enabled",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
			},
		},
		{
			name: "Failed on GetWalletByCustomerXid",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, mockError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewWalletUsecase(mockWalletDB)
			_, err := uc.Enabled(tt.args.ctx, tt.args.customerXid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Enabled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseWallet_Disable(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletDB := mockDB.NewMockWalletRepository(ctrl)
	defer ctrl.Finish()
	now := time.Now()
	mockCtx := context.Background()
	mockError := errors.New("fake error")
	mockWallet := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockWalletDisable := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	mockCustomerXid := mockWallet.OwnedBy

	type args struct {
		ctx         context.Context
		customerXid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "Success",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: false,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockWalletDB.EXPECT().UpdateStatusWallet(gomock.Any(), gomock.Any()).Return(mockWallet, nil)
			},
		},
		{
			name: "Failed on GetWalletByCustomerXid",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, mockError)
			},
		},
		{
			name: "Failed wallet empty",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(nil, nil)
			},
		},
		{
			name: "Failed wallet already disabled",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWalletDisable, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewWalletUsecase(mockWalletDB)
			_, err := uc.Disable(tt.args.ctx, tt.args.customerXid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Disable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
