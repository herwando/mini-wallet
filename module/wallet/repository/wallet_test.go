package repository_test

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/repository"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	now          = time.Now()
	walletEnable = &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    1,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(10000),
	}
	walletDisable = &model.Wallet{
		ID:        "81401b03-60e0-4f20-afc6-419b3773e7b3",
		OwnedBy:   "ea0212d3-abd6-406f-8c67-868e814a2436",
		Status:    2,
		EnabledAt: now,
		Balance:   decimal.NewFromInt(0),
	}
)

func TestWallet_CreateWallet(t *testing.T) {
	testCases := map[string]struct {
		req     *model.Wallet
		result  *model.Wallet
		wantErr bool
	}{
		"success": {
			req:     walletEnable,
			result:  walletEnable,
			wantErr: false,
		},
		"failed scan": {
			req:     walletEnable,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewWalletRepository(db)

			column := []string{"id"}
			row := sMock.NewRows(column).AddRow(walletEnable.ID)
			if tc.wantErr {
				column = []string{"id", "id"}
				row = sMock.NewRows(column).AddRow(walletEnable.ID, walletEnable.ID)
			}

			sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO wallets (owned_by, status, enabled_at, balance) VALUES ($1, $2, $3, $4) RETURNING id")).
				WithArgs(walletEnable.OwnedBy, walletEnable.Status, walletEnable.EnabledAt, walletEnable.Balance).WillReturnRows(row)

			result, err := repo.CreateWallet(context.Background(), tc.req)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.CreateWallet() = %v, want %v", result, tc.result)
			}
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWallet_UpdateStatusWallet(t *testing.T) {
	testCases := map[string]struct {
		req     *model.Wallet
		result  *model.Wallet
		wantErr bool
	}{
		"success": {
			req:     walletDisable,
			result:  walletDisable,
			wantErr: false,
		},
		"failed scan": {
			req:     walletEnable,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewWalletRepository(db)

			column := []string{"id"}
			row := sMock.NewRows(column).AddRow(walletDisable.ID)
			if tc.wantErr {
				column = []string{"id", "id"}
				row = sMock.NewRows(column).AddRow(walletDisable.ID, walletDisable.ID)
			}

			sMock.ExpectQuery(regexp.QuoteMeta("UPDATE wallets SET status = $1 WHERE id = $2 RETURNING id")).
				WithArgs(walletDisable.Status, walletDisable.ID).WillReturnRows(row)

			result, err := repo.UpdateStatusWallet(context.Background(), tc.req)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.UpdateStatusWallet() = %v, want %v", result, tc.result)
			}
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWallet_GetWalletByCustomerXid(t *testing.T) {
	testCases := map[string]struct {
		req     string
		result  *model.Wallet
		wantErr bool
		err     error
	}{
		"success": {
			req:     walletEnable.OwnedBy,
			result:  walletEnable,
			wantErr: false,
		},
		"failed scan": {
			req:     walletEnable.OwnedBy,
			wantErr: true,
		},
		"failed sql no row": {
			req:     walletEnable.OwnedBy,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewWalletRepository(db)

			column := []string{"id", "owned_by", "status", "enabled_at", "balance"}
			row := sMock.NewRows(column).AddRow(walletEnable.ID, walletEnable.OwnedBy, walletEnable.Status, walletEnable.EnabledAt, walletEnable.Balance)

			if tc.err != nil {
				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, owned_by, status, enabled_at, balance FROM wallets WHERE owned_by = $1")).
					WithArgs(tc.req).WillReturnError(tc.err)
			} else {
				if tc.wantErr {
					column = []string{"id", "id"}
					row = sMock.NewRows(column).AddRow(walletEnable.ID, walletEnable.ID)
				}

				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id, owned_by, status, enabled_at, balance FROM wallets WHERE owned_by = $1")).
					WithArgs(tc.req).WillReturnRows(row)
			}

			result, _ := repo.GetWalletByCustomerXid(context.Background(), tc.req)

			if !reflect.DeepEqual(result, tc.result) {
				t.Errorf("Wallet.GetWalletByCustomerXid() = %v, want %v", result, tc.result)
			}
		})
	}
}
