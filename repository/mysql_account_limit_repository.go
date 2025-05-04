package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type accountLimitRepositoryMysql struct {
	dbtx DBTX
}

func NewAccountLimitRepositoryMysql(dbtx DBTX) *accountLimitRepositoryMysql {
	return &accountLimitRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *accountLimitRepositoryMysql) GetAccountLimitByAccountId(ctx context.Context, accountId int64, forUpdate bool) (*entity.AccountLimit, error) {
	var sb strings.Builder

	sb.WriteString(`
		SELECT 
			account_limit_id,
			account_id,
			account_limit_1_m,
			account_limit_2_m,
			account_limit_3_m,
			account_limit_4_m,
			created_at,
			updated_at,
			deleted_at
		FROM account_limits
		WHERE account_id = ?
			AND deleted_at IS NULL
	`)

	if forUpdate {
		sb.WriteString(`FOR UPDATE`)
	}

	q := sb.String()

	var accountLimit entity.AccountLimit

	err := r.dbtx.QueryRowContext(ctx, q, accountId).Scan(
		&accountLimit.Id,
		&accountLimit.AccountId,
		&accountLimit.Limit1M,
		&accountLimit.Limit2M,
		&accountLimit.Limit3M,
		&accountLimit.Limit4M,
		&accountLimit.CreatedAt,
		&accountLimit.UpdatedAt,
		&accountLimit.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &accountLimit, nil
}

