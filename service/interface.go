package service

import (
	"context"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type AccountService interface {
	RegisterAccount(ctx context.Context, newAccount entity.Account) (*entity.TokenData, error)
}
