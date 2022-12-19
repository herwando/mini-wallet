package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/repository"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	deposit = &model.Deposit{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		DepositedBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		DepositedAt: now,
		Amount:      decimal.NewFromInt(1000),
		ReferenceID: "81401b03-60e0-4f20-afc6-419b3773e7b3",
	}
)

func TestDeposit_GetDepositByReferenceId(t *testing.T) {
	testCases := map[string]struct {
		req     string
		result  *model.Deposit
		wantErr bool
		err     error
	}{
		"success": {
			req:     deposit.ReferenceID,
			result:  deposit,
			wantErr: false,
		},
		"failed scan": {
			req:     deposit.ReferenceID,
			wantErr: true,
		},
		"failed sql no row": {
			req:     deposit.ReferenceID,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewDepositRepository(db)

			column := []string{"id", "deposited_by", "status", "deposited_at", "amount", "reference_id"}
			row := sMock.NewRows(column).AddRow(deposit.ID, deposit.DepositedBy, deposit.Status, deposit.DepositedAt, deposit.Amount, deposit.ReferenceID)

			if tc.err != nil {
				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, deposited_by, status, deposited_at, amount, reference_id FROM deposits WHERE reference_id = $1")).
					WithArgs(tc.req).WillReturnError(tc.err)
			} else {
				if tc.wantErr {
					column = []string{"id", "id"}
					row = sMock.NewRows(column).AddRow(deposit.ID, deposit.ID)
				}

				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, deposited_by, status, deposited_at, amount, reference_id FROM deposits WHERE reference_id = $1")).
					WithArgs(tc.req).WillReturnRows(row)
			}

			result, _ := repo.GetDepositByReferenceId(context.Background(), tc.req)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.GetDepositByReferenceId() = %v, want %v", result, tc.result)
			}
		})
	}
}

func TestDeposit_CreateDeposit(t *testing.T) {
	testCases := map[string]struct {
		reqDeposit     *model.Deposit
		reqWallet      *model.Wallet
		result         *model.Deposit
		wantErr        bool
		wantErrTx      bool
		wantErrDeposit bool
		wantErrWallet  bool
		wantErrCommit  bool
		err            error
	}{
		"success": {
			reqDeposit:     deposit,
			reqWallet:      walletEnable,
			result:         deposit,
			wantErr:        false,
			wantErrTx:      false,
			wantErrDeposit: false,
			wantErrWallet:  false,
			wantErrCommit:  false,
		},
		"failed tx": {
			reqDeposit: deposit,
			reqWallet:  walletEnable,
			wantErr:    true,
			wantErrTx:  true,
			err:        errors.New("database error"),
		},
		"failed deposit": {
			reqDeposit:     deposit,
			reqWallet:      walletEnable,
			wantErr:        true,
			wantErrTx:      false,
			wantErrDeposit: true,
			err:            errors.New("database error"),
		},
		"failed wallet": {
			reqDeposit:     deposit,
			reqWallet:      walletEnable,
			wantErr:        true,
			wantErrTx:      false,
			wantErrDeposit: false,
			wantErrWallet:  true,
			err:            errors.New("database error"),
		},
		"failed commit": {
			reqDeposit:     deposit,
			reqWallet:      walletEnable,
			wantErr:        true,
			wantErrTx:      false,
			wantErrDeposit: false,
			wantErrWallet:  false,
			wantErrCommit:  true,
			err:            errors.New("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewDepositRepository(db)

			if tc.wantErrTx {
				sMock.ExpectBegin().WillReturnError(tc.err)
				sMock.ExpectRollback()
			} else {
				sMock.ExpectBegin()
				column := []string{"id"}
				row := sMock.NewRows(column).AddRow(deposit.ID)
				if tc.wantErrDeposit {
					sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO deposits (deposited_by, status, deposited_at, amount, reference_id) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
						WithArgs(tc.reqDeposit.DepositedBy, tc.reqDeposit.Status, tc.reqDeposit.DepositedAt, tc.reqDeposit.Amount, tc.reqDeposit.ReferenceID).WillReturnError(tc.err)
					sMock.ExpectRollback()
				} else {
					sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO deposits (deposited_by, status, deposited_at, amount, reference_id) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
						WithArgs(tc.reqDeposit.DepositedBy, tc.reqDeposit.Status, tc.reqDeposit.DepositedAt, tc.reqDeposit.Amount, tc.reqDeposit.ReferenceID).WillReturnRows(row)
					row = sMock.NewRows(column).AddRow(walletEnable.ID)
					if tc.wantErrWallet {
						sMock.ExpectQuery(regexp.QuoteMeta("UPDATE wallets SET balance = $1 WHERE id = $2 RETURNING id")).
							WithArgs(tc.reqWallet.Balance, tc.reqWallet.ID).WillReturnError(tc.err)
						sMock.ExpectRollback()
					} else {
						sMock.ExpectQuery(regexp.QuoteMeta("UPDATE wallets SET balance = $1 WHERE id = $2 RETURNING id")).
							WithArgs(tc.reqWallet.Balance, tc.reqWallet.ID).WillReturnRows(row)
						if tc.wantErrCommit {
							sMock.ExpectCommit().WillReturnError(tc.err)
							sMock.ExpectRollback()
						} else {
							sMock.ExpectCommit()
						}
					}
				}
			}

			result, err := repo.CreateDeposit(context.Background(), tc.reqDeposit, tc.reqWallet)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.CreateDeposit() = %v, want %v", result, tc.result)
			}

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
