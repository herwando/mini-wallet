package repository

import (
	"context"
	"database/sql"

	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO
		wallets (owned_by, status, enabled_at, balance) 
		VALUES ($1, $2, $3, $4) RETURNING id`, wallet.OwnedBy, wallet.Status, wallet.EnabledAt, wallet.Balance)
	err := row.Scan(&wallet.ID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) UpdateStatusWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	row := r.db.QueryRowContext(ctx, `UPDATE wallets
		SET status = $1 WHERE id = $2 RETURNING id`, wallet.Status, wallet.ID)
	err := row.Scan(&wallet.ID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletRepository) GetWalletByCustomerXid(ctx context.Context, customerXid string) (*model.Wallet, error) {
	query := `
		SELECT
			id, owned_by, status, enabled_at, balance
		FROM wallets WHERE owned_by = $1
	`

	wallet := &model.Wallet{}
	row := r.db.QueryRowContext(ctx, query, customerXid)
	if err := row.Scan(&wallet.ID, &wallet.OwnedBy, &wallet.Status, &wallet.EnabledAt, &wallet.Balance); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return wallet, nil
}
