package repository

import (
	"context"
	"fmt"
	"strings"
)

type refreshTokenRepositoryMysql struct {
	dbtx DBTX
}

func NewRefreshTokenRepositoryMysql(dbtx DBTX) *refreshTokenRepositoryMysql {
	return &refreshTokenRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *refreshTokenRepositoryMysql) InsertToken(ctx context.Context, token string, accountId, expiredAt int64) error {
	var sb strings.Builder
	
	sb.WriteString(`
		INSERT INTO refresh_tokens (refresh_token, account_id, expired_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`)

	q := sb.String()

	now := nowUnixMilli()

	_, err := r.dbtx.ExecContext(ctx, q, token, accountId, expiredAt, now, now)
	if err != nil {
		return fmt.Errorf("[mysql_refresh_token_repository][InsertToken][ExecContext] error: %w", err)
	}

	return nil
}
