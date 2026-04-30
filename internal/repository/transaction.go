package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"hechi-go/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*domain.Transaction, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO transactions (id, account_id, amount, type, category, description, date, created_at)
		 VALUES ($1, $2, $3, $4::transaction_type, $5, $6, $7, $8)`,
		tx.ID, tx.AccountID, tx.Amount.String(), string(tx.Type),
		tx.Category, tx.Description, tx.Date, tx.CreatedAt,
	)
	return err
}

func (r *transactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*domain.Transaction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, account_id, amount::text, type::text, category, description, date, created_at
		 FROM transactions WHERE account_id = $1
		 ORDER BY date DESC LIMIT $2 OFFSET $3`,
		accountID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		var amountStr, txType string
		if err := rows.Scan(&tx.ID, &tx.AccountID, &amountStr, &txType,
			&tx.Category, &tx.Description, &tx.Date, &tx.CreatedAt); err != nil {
			return nil, err
		}
		tx.Type = domain.TransactionType(txType)
		tx.Amount, _ = decimal.NewFromString(amountStr)
		txs = append(txs, &tx)
	}
	return txs, rows.Err()
}

func (r *transactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	return err
}
