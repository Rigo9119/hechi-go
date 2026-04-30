package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeCredit     AccountType = "CREDIT"
	AccountTypeInvestment AccountType = "INVESTMENT"
)

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      AccountType
	Balance   decimal.Decimal
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
