package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
	mockDB "github.com/herwando/mini-wallet/module/wallet/usecase/_mocks"
)

func TestUsecaseAccount_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAccountDB := mockDB.NewMockAccountRepository(ctrl)
	defer ctrl.Finish()

	mockCtx := context.Background()
	mockError := errors.New("fake error")
	mockAccount := model.Account{
		CustomerXid: "ea0212d3-abd6-406f-8c67-868e814a2436",
	}

	type args struct {
		ctx     context.Context
		payload model.Account
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
				ctx:     mockCtx,
				payload: mockAccount,
			},
			wantErr: false,
			patch: func() {
				mockAccountDB.EXPECT().ExistAccountByCustomerXid(gomock.Any(), mockAccount.CustomerXid).Return(false, nil)
				mockAccountDB.EXPECT().CreateAccount(gomock.Any(), &mockAccount).Return(nil)
			},
		},
		{
			name: "Success without create",
			args: args{
				ctx:     mockCtx,
				payload: mockAccount,
			},
			wantErr: false,
			patch: func() {
				mockAccountDB.EXPECT().ExistAccountByCustomerXid(gomock.Any(), mockAccount.CustomerXid).Return(true, nil)
			},
		},
		{
			name: "Failed on ExistAccountByCustomerXid cause databaseErr",
			args: args{
				ctx:     mockCtx,
				payload: mockAccount,
			},
			wantErr: true,
			patch: func() {
				mockAccountDB.EXPECT().ExistAccountByCustomerXid(gomock.Any(), mockAccount.CustomerXid).Return(false, mockError)
			},
		},
		{
			name: "Failed on CreateAccount cause databaseErr",
			args: args{
				ctx:     mockCtx,
				payload: mockAccount,
			},
			wantErr: true,
			patch: func() {
				mockAccountDB.EXPECT().ExistAccountByCustomerXid(gomock.Any(), mockAccount.CustomerXid).Return(false, nil)
				mockAccountDB.EXPECT().CreateAccount(gomock.Any(), &mockAccount).Return(mockError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.patch()
			uc := usecase.NewAccountUsecase(mockAccountDB)
			_, err := uc.Init(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
