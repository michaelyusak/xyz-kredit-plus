package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *accountRepositoryMysql) GetAccountByEmail(ctx context.Context, email string, forUpdate bool) (*entity.Account, error) {
	var sb strings.Builder

	sb.WriteString(`
		SELECT account_id, email, password, created_at, updated_at, deleted_at
		FROM accounts
		WHERE email = ?
			AND deleted_at IS NULL
	`)

	if forUpdate {
		sb.WriteString(`FOR UPDATE`)
	}

	q := sb.String()

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
	var sb strings.Builder
	
	sb.WriteString(`
		INSERT INTO accounts (email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`)

	q := sb.String()
	
	now := nowUnixMilli()

	res, err := r.dbtx.ExecContext(ctx, q, account.Email, account.Password, now, now)
	if err != nil {
		return 0, fmt.Errorf("[mysql_account_repository][InsertAccount][ExecContext] error: %w | email: %s", err, account.Email)
	}

	accountId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[mysql_account_repository][InsertAccount][LastInsertId] error: %w | email: %s", err, account.Email)
	}

	return accountId, nil
}
