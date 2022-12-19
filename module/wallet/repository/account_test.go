package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
	"github.com/herwando/mini-wallet/module/wallet/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var (
	account = &model.Account{
		CustomerXid: "ea0212d3-abd6-406f-8c67-868e814a2436",
	}
)

func TestAccount_CreateAccount(t *testing.T) {
	testCases := map[string]struct {
		req     *model.Account
		wantErr bool
	}{
		"success": {
			req:     account,
			wantErr: false,
		},
		"failed scan": {
			req:     account,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewAccountRepository(db)

			column := []string{"id"}
			row := sMock.NewRows(column).AddRow(account.CustomerXid)
			if tc.wantErr {
				column = []string{"id", "id"}
				row = sMock.NewRows(column).AddRow(account.CustomerXid, account.CustomerXid)
			}

			sMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (id) VALUES ($1)")).
				WithArgs(account.CustomerXid).WillReturnRows(row)

			err := repo.CreateAccount(context.Background(), tc.req)

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAccount_ExistAccountByCustomerXid(t *testing.T) {
	testCases := map[string]struct {
		req     string
		want    bool
		wantErr bool
		err     error
	}{
		"success": {
			req:     account.CustomerXid,
			want:    true,
			wantErr: false,
		},
		"failed scan": {
			req:     account.CustomerXid,
			want:    false,
			wantErr: true,
		},
		"failed sql no row": {
			req:     account.CustomerXid,
			want:    false,
			wantErr: true,
			err:     sql.ErrNoRows,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockDB, sMock, _ := sqlmock.New()
			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
			db := sqlxDB.DB
			repo := repository.NewAccountRepository(db)

			column := []string{"id"}
			row := sMock.NewRows(column).AddRow(account.CustomerXid)
			if tc.err != nil {
				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM accounts WHERE id = $1")).
					WithArgs(tc.req).WillReturnError(tc.err)
			} else {
				if tc.wantErr {
					column = []string{"id", "id"}
					row = sMock.NewRows(column).AddRow(account.CustomerXid, account.CustomerXid)
				}

				sMock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM accounts WHERE id = $1")).
					WithArgs(tc.req).WillReturnRows(row)
			}

			result, _ := repo.ExistAccountByCustomerXid(context.Background(), tc.req)
			assert.Equal(t, tc.want, result)
		})
	}
}
