package graph

import (
	"hechi-go/internal/auth"
	"hechi-go/internal/repository"
)

type Resolver struct {
	UserRepo        repository.UserRepository
	AccountRepo     repository.AccountRepository
	TransactionRepo repository.TransactionRepository
	Auth            *auth.Service
}
