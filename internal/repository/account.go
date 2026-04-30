package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"hechi-go/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance decimal.Decimal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, a *domain.Account) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO accounts (id, user_id, name, type, balance, currency, created_at, updated_at)
		 VALUES ($1, $2, $3, $4::account_type, $5, $6, $7, $8)`,
		a.ID, a.UserID, a.Name, string(a.Type), a.Balance.String(), a.Currency, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *accountRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var a domain.Account
	var balanceStr, accountType string
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, type::text, balance::text, currency, created_at, updated_at
		 FROM accounts WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.UserID, &a.Name, &accountType, &balanceStr, &a.Currency, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	a.Type = domain.AccountType(accountType)
	a.Balance, _ = decimal.NewFromString(balanceStr)
	return &a, nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, type::text, balance::text, currency, created_at, updated_at
		 FROM accounts WHERE user_id = $1 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var a domain.Account
		var balanceStr, accountType string
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &accountType, &balanceStr, &a.Currency, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.Type = domain.AccountType(accountType)
		a.Balance, _ = decimal.NewFromString(balanceStr)
		accounts = append(accounts, &a)
	}
	return accounts, rows.Err()
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id uuid.UUID, balance decimal.Decimal) error {
	_, err := r.db.Exec(ctx,
		`UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`,
		balance.String(), id,
	)
	return err
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM accounts WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return err
}
