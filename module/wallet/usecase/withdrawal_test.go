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

func TestUsecaseWithdrawal_CreateWithdrawal(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWithdrawalDB := mockDB.NewMockWithdrawalRepository(ctrl)
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
	mockWalletZero := &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(0),
	}
	mockWithdrawal := &model.Withdrawal{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		WithdrawnBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		WithdrawnAt: now,
		Amount:      decimal.NewFromInt(1000),
		ReferenceID: "81401b03-60e0-4f20-afc6-419b3773e7b3",
	}
	mockCustomerXid := mockWithdrawal.WithdrawnBy
	mockPayload := model.PayloadWithdrawal{
		ReferenceID: mockWithdrawal.ReferenceID,
		Amount:      mockWithdrawal.Amount,
	}

	type args struct {
		ctx         context.Context
		customerXid string
		payload     model.PayloadWithdrawal
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
				mockWithdrawalDB.EXPECT().GetWithdrawalByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, nil)
				mockWithdrawalDB.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockWithdrawal, nil)
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
			name: "Failed wallet less than balance",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWalletZero, nil)
			},
		},
		{
			name: "Failed on GetWithdrawalByReferenceId errpr",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockWithdrawalDB.EXPECT().GetWithdrawalByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, mockError)
			},
		},
		{
			name: "Failed on GetWithdrawalByReferenceId",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockWithdrawalDB.EXPECT().GetWithdrawalByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(mockWithdrawal, nil)
			},
		},
		{
			name: "Failed on CreateWithdrawal cause databaseErr",
			args: args{
				ctx:         mockCtx,
				customerXid: mockCustomerXid,
				payload:     mockPayload,
			},
			wantErr: true,
			patch: func() {
				mockWalletDB.EXPECT().GetWalletByCustomerXid(gomock.Any(), mockCustomerXid).Return(mockWallet, nil)
				mockWithdrawalDB.EXPECT().GetWithdrawalByReferenceId(gomock.Any(), mockPayload.ReferenceID).Return(nil, nil)
				mockWithdrawalDB.EXPECT().CreateWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, mockError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewWithdrawalUsecase(mockWithdrawalDB, mockWalletDB)
			_, err := uc.CreateWithdrawal(tt.args.ctx, tt.args.customerXid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.CreateWithdrawal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
