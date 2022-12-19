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

func TestUsecaseDeposit_CreateDeposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDepositDB := mockDB.NewMockDepositRepository(ctrl)
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
	mockDeposit := &model.Deposit{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		DepositedBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		DepositedAt: now,
		Amount:      decimal.NewFromInt(1000),
		ReferenceID: "81401b03-60e0-4f20-afc6-419b3773e7b3",
	}
	mockCustomerXid := mockDeposit.DepositedBy
	mockPayload := model.PayloadDeposit{
		ReferenceID: mockDeposit.ReferenceID,
		Amount:      mockDeposit.Amount,
	}

	type args struct {
		ctx         context.Context
		customerXid string
		payload     model.PayloadDeposit
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
				payload:     mockPayload,
			},
			wantErr: false,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockDepositDB.EXPECT().GetDepositByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, nil)
				mockDepositDB.EXPECT().CreateDeposit(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDeposit, nil)
			},
		},
		{
			name: "Failed on GetWalletByCustomerXid errpr",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
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
				payload:     mockPayload,
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
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWalletDisable, nil)
			},
		},
		{
			name: "Failed on GetDepositByReferenceId errpr",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockDepositDB.EXPECT().GetDepositByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, mockError)
			},
		},
		{
			name: "Failed on GetDepositByReferenceId",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockDepositDB.EXPECT().GetDepositByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(mockDeposit, nil)
			},
		},
		{
			name: "Failed on CreateDeposit cause databaseErr",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockDepositDB.EXPECT().GetDepositByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, nil)
				mockDepositDB.EXPECT().CreateDeposit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mockError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewDepositUsecase(mockDepositDB, mockWalletDB)
			_, err := uc.CreateDeposit(tt.args.ctx, tt.args.customerXid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.CreateDeposit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
