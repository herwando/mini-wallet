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
	withdrawal = &model.Withdrawal{
		ID:          "81401b03-60e0-4f20-afc6-419b3773e7b3",
		WithdrawnBy: "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:      1,
		WithdrawnAt: now,
		Amount:      decimal.NewFromInt(1000),
		ReferenceID: "81401b03-60e0-4f20-afc6-419b3773e7b3",
	}
)

func TestWithdrawal_GetWithdrawalByReferenceId(t *testing.T) {
	testCases := map[string]struct {
		req     string
		result  *model.Withdrawal
		wantErr bool
		err     error
	}{
		"success": {
			req:     withdrawal.ReferenceID,
			result:  withdrawal,
			wantErr: false,
		},
		"failed scan": {
			req:     withdrawal.ReferenceID,
			wantErr: true,
		},
		"failed sql no row": {
			req:     withdrawal.ReferenceID,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewWithdrawalRepository(db)

			column := []string{"id", "withdrawn_by", "status", "withdrawn_at", "amount", "reference_id"}
			row := sMock.NewRows(column).AddRow(withdrawal.ID, withdrawal.WithdrawnBy, withdrawal.Status, withdrawal.WithdrawnAt, withdrawal.Amount, withdrawal.ReferenceID)

			if tc.err != nil {
				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, withdrawn_by, status, withdrawn_at, amount, reference_id FROM withdrawals WHERE reference_id = $1")).
					WithArgs(tc.req).WillReturnError(tc.err)
			} else {
				if tc.wantErr {
					column = []string{"id", "id"}
					row = sMock.NewRows(column).AddRow(deposit.ID, deposit.ID)
				}

				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, withdrawn_by, status, withdrawn_at, amount, reference_id FROM withdrawals WHERE reference_id = $1")).
					WithArgs(tc.req).WillReturnRows(row)
			}

			result, _ := repo.GetWithdrawalByReferenceId(context.Background(), tc.req)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.GetWithdrawalByReferenceId() = %v, want %v", result, tc.result)
			}
		})
	}
}

func TestWithdrawal_CreateWithdrawal(t *testing.T) {
	testCases := map[string]struct {
		reqWithdrawal     *model.Withdrawal
		reqWallet         *model.Wallet
		result            *model.Withdrawal
		wantErr           bool
		wantErrTx         bool
		wantErrWithdrawal bool
		wantErrWallet     bool
		wantErrCommit     bool
		err               error
	}{
		"success": {
			reqWithdrawal:     withdrawal,
			reqWallet:         walletEnable,
			result:            withdrawal,
			wantErr:           false,
			wantErrTx:         false,
			wantErrWithdrawal: false,
			wantErrWallet:     false,
			wantErrCommit:     false,
		},
		"failed tx": {
			reqWithdrawal: withdrawal,
			reqWallet:     walletEnable,
			wantErr:       true,
			wantErrTx:     true,
			err:           errors.New("database error"),
		},
		"failed withdrawal": {
			reqWithdrawal:     withdrawal,
			reqWallet:         walletEnable,
			wantErr:           true,
			wantErrTx:         false,
			wantErrWithdrawal: true,
			err:               errors.New("database error"),
		},
		"failed wallet": {
			reqWithdrawal:     withdrawal,
			reqWallet:         walletEnable,
			wantErr:           true,
			wantErrTx:         false,
			wantErrWithdrawal: false,
			wantErrWallet:     true,
			err:               errors.New("database error"),
		},
		"failed commit": {
			reqWithdrawal:     withdrawal,
			reqWallet:         walletEnable,
			wantErr:           true,
			wantErrTx:         false,
			wantErrWithdrawal: false,
			wantErrWallet:     false,
			wantErrCommit:     true,
			err:               errors.New("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewWithdrawalRepository(db)

			if tc.wantErrTx {
				sMock.ExpectBegin().WillReturnError(tc.err)
				sMock.ExpectRollback()
			} else {
				sMock.ExpectBegin()
				column := []string{"id"}
				row := sMock.NewRows(column).AddRow(withdrawal.ID)
				if tc.wantErrWithdrawal {
					sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO withdrawals (withdrawn_by, status, withdrawn_at, amount, reference_id) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
						WithArgs(tc.reqWithdrawal.WithdrawnBy, tc.reqWithdrawal.Status, tc.reqWithdrawal.WithdrawnAt, tc.reqWithdrawal.Amount, tc.reqWithdrawal.ReferenceID).WillReturnError(tc.err)
					sMock.ExpectRollback()
				} else {
					sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO withdrawals (withdrawn_by, status, withdrawn_at, amount, reference_id) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
						WithArgs(tc.reqWithdrawal.WithdrawnBy, tc.reqWithdrawal.Status, tc.reqWithdrawal.WithdrawnAt, tc.reqWithdrawal.Amount, tc.reqWithdrawal.ReferenceID).WillReturnRows(row)
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

			result, err := repo.CreateWithdrawal(context.Background(), tc.reqWithdrawal, tc.reqWallet)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.CreateWithdrawal() = %v, want %v", result, tc.result)
			}

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
