package repository

import (
	"context"
	"database/sql"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WithdrawalRepository struct {
	db *sql.DB
}

type WithdrawalRepositoryInterface interface {
	GetWithdrawalByReferenceId(ctx context.Context, referenceId string) (*model.Withdrawal, error)
	CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal, wallet *model.Wallet) (*model.Withdrawal, error)
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

func (r *WithdrawalRepository) GetWithdrawalByReferenceId(ctx context.Context, referenceId string) (*model.Withdrawal, error) {
	query := `
		SELECT
			id, withdrawn_by, status, withdrawn_at, amount, reference_id
		FROM withdrawals WHERE reference_id = $1
	`

	withdrawal := &model.Withdrawal{}
	row := r.db.QueryRowContext(ctx, query, referenceId)
	if err := row.Scan(&withdrawal.ID, &withdrawal.WithdrawnBy, &withdrawal.Status, &withdrawal.WithdrawnAt, &withdrawal.Amount, &withdrawal.ReferenceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return withdrawal, nil
}

func (r *WithdrawalRepository) CreateWithdrawal(ctx context.Context, withdrawal *model.Withdrawal, wallet *model.Wallet) (*model.Withdrawal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO
		withdrawals (withdrawn_by, status, withdrawn_at, amount, reference_id) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, withdrawal.WithdrawnBy, withdrawal.Status, withdrawal.WithdrawnAt, withdrawal.Amount, withdrawal.ReferenceID)
	err = row.Scan(&withdrawal.ID)
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

	return withdrawal, nil
}
