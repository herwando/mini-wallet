package repository

import (
	"context"
	"database/sql"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type DepositRepository struct {
	db *sql.DB
}

func NewDepositRepository(db *sql.DB) *DepositRepository {
	return &DepositRepository{
		db: db,
	}
}

func (r *DepositRepository) GetDepositByReferenceId(ctx context.Context, referenceId string) (*model.Deposit, error) {
	query := `
		SELECT
			id, deposited_by, status, deposited_at, amount, reference_id
		FROM deposits WHERE reference_id = $1
	`

	deposit := &model.Deposit{}
	row := r.db.QueryRowContext(ctx, query, referenceId)
	if err := row.Scan(&deposit.ID, &deposit.DepositedBy, &deposit.Status, &deposit.DepositedAt, &deposit.Amount, &deposit.ReferenceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return deposit, nil
}

func (r *DepositRepository) CreateDeposit(ctx context.Context, deposit *model.Deposit, wallet *model.Wallet) (*model.Deposit, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO
		deposits (deposited_by, status, deposited_at, amount, reference_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, deposit.DepositedBy, deposit.Status, deposit.DepositedAt, deposit.Amount, deposit.ReferenceID)
	err = row.Scan(&deposit.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	row = tx.QueryRowContext(ctx, `UPDATE wallets
		SET balance = $1 WHERE id = $2 RETURNING id`, wallet.Balance, wallet.ID)
	err = row.Scan(&wallet.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return deposit, nil
}
