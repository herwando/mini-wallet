package repository

import (
	"context"
	"database/sql"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) ExistAccountByCustomerXid(ctx context.Context, customerXid string) (bool, error) {
	query := `
		SELECT
			id
		FROM accounts WHERE id = $1
	`

	account := model.Account{}
	row := r.db.QueryRowContext(ctx, query, customerXid)
	if err := row.Scan(&account); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return true, err
	}

	return false, nil
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO
		accounts (id) 
		VALUES (:id)`, &account)
	if err != nil {
		return err
	}

	return nil
}
