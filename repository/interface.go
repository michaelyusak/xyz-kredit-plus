package repository

import (
	"context"
	"database/sql"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type AccountRepository interface {
	GetAccountByEmail(ctx context.Context, email string, forUpdate bool) (*entity.Account, error)
	InsertAccount(ctx context.Context, account entity.Account) (int64, error)
}

type RefreshTokenRepository interface {
	InsertToken(ctx context.Context, token string, accountId, expiredAt int64) error
}

type ConsumerRepository interface {
	GetConsumerByAccountId(ctx context.Context, accountId int64, forUpdate bool) (*entity.Consumer, error)
}

type MediaRepository interface {
	Store(ctx context.Context, media MediaOpt) error
}

type AccountLimitRepository interface {
	GetAccountLimitByAccountId(ctx context.Context, accountId int64, forUpdate bool) (*entity.AccountLimit, error)
	UpdateLimit(ctx context.Context, accountLimit entity.AccountLimit) error
}

type TransactionRepository interface {
	InsertTransaction(ctx context.Context, transaction entity.Transaction) (int64, error)
}
