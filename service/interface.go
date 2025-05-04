package service

import (
	"context"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type AccountService interface {
	RegisterAccount(ctx context.Context, newAccount entity.Account) (*entity.TokenData, error)
	Login(ctx context.Context, account entity.Account) (*entity.TokenData, error)
}

type ConsumerService interface {
	ProcessKyc(ctx context.Context, consumerData entity.Consumer) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error)
}
