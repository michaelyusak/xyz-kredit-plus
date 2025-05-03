package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type accountRepositoryMysql struct {
	dbtx DBTX
}

func NewAccountRepositoryMysql(dbtx DBTX) *accountRepositoryMysql {
	return &accountRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *accountRepositoryMysql) Lock(ctx context.Context) error {
	q := `
		LOCK TABLE accounts WRITE
	`

	_, err := r.dbtx.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("[mysql_account_repository][Lock][ExecContext] error: %w", err)
	}

	return nil
}

func (r *accountRepositoryMysql) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	q := `
		SELECT account_id, email, password, created_at, updated_at, deleted_at
		WHERE email = $1
			AND deleted_at IS NULL
	`

	var account entity.Account

	err := r.dbtx.QueryRowContext(ctx, q, email).Scan(
		&account.Id,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("[mysql_account_repository][GetAccountByEmail][QueryRowContext] error: %w | email: %s", err, email)
	}

	return &account, nil
}

func (r *accountRepositoryMysql) InsertAccount(ctx context.Context, account entity.Account) (int64, error) {
	q := `
		INSERT INTO accounts (email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $3)
	`

	res, err := r.dbtx.ExecContext(ctx, q, account.Email, account.Password, nowUnixMilli())
	if err != nil {
		return 0, fmt.Errorf("[mysql_account_repository][InsertAccount][ExecContext] error: %w | email: %s", err, account.Email)
	}

	accountId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[mysql_account_repository][InsertAccount][LastInsertId] error: %w | email: %s", err, account.Email)
	}

	return accountId, nil
}
