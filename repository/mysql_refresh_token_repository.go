package repository

import (
	"context"
	"fmt"
)

type refreshTokenRepositoryMysql struct {
	dbtx DBTX
}

func NewRefreshTokenRepositoryMysql(dbtx DBTX) *refreshTokenRepositoryMysql {
	return &refreshTokenRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *refreshTokenRepositoryMysql) Lock(ctx context.Context) error {
	q := `
		LOCK TABLE refresh_tokens WRITE
	`

	_, err := r.dbtx.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("[mysql_refresh_token_repository][Lock][ExecContext] error: %w", err)
	}

	return nil
}

func (r *refreshTokenRepositoryMysql) InsertToken(ctx context.Context, token string, accountId int64) error {
	q := `
	INSERT INTO refresh_tokens (refresh_token, account_id, created_at, updated_at)
	VALUES ($1, $2, $3, $3)
`

	_, err := r.dbtx.ExecContext(ctx, q, token, accountId, nowUnixMilli())
	if err != nil {
		return fmt.Errorf("[mysql_refresh_token_repository][InsertToken][ExecContext] error: %w", err)
	}

	return nil
}
