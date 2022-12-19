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
	if err := row.Scan(&account.CustomerXid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	row := r.db.QueryRowContext(ctx, `INSERT INTO
		accounts (id) 
		VALUES ($1) RETURNING id`, account.CustomerXid)
	err := row.Scan(&account.CustomerXid)
	if err != nil {
		return err
	}

	return nil
}
