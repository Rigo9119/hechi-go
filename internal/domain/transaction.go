package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "INCOME"
	TransactionTypeExpense  TransactionType = "EXPENSE"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Amount      decimal.Decimal
	Type        TransactionType
	Category    *string
	Description *string
	Date        time.Time
	CreatedAt   time.Time
}
