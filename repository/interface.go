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
	Lock(ctx context.Context) error
	GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error)
	InsertAccount(ctx context.Context, account entity.Account) (error, int64)
}

type RefreshTokenRepository interface {
	Lock(ctx context.Context) error
	InsertToken(ctx context.Context, token string, accountId int64) error
}

type ConsumerRepository interface {
	Lock(ctx context.Context) error
	GetConsumetByAccountId(ctx context.Context, accountId int64) (*entity.Consumer, error)
}
