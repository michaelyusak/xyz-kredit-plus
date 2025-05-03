package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type consumerRepositoryMysql struct {
	dbtx DBTX
}

func NewConsumerRepositoryMysql(dbtx DBTX) *consumerRepositoryMysql {
	return &consumerRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *consumerRepositoryMysql) GetConsumerByAccountId(ctx context.Context, accountId int64, forUpdate bool) (*entity.Consumer, error) {
	var sb strings.Builder

	sb.WriteString(`
		SELECT 
			consumer_id, 
			account_id, 
			identity_number, 
			full_name, 
			legal_name, 
			place_of_birth, 
			date_of_birth, 
			salary, 
			identity_card_photo_url, 
			selfie_photo_url, 
			created_at, 
			updated_at, 
			deleted_at
		FROM consumers
		WHERE account_id = ?
		 	AND deleted_at IS NULL
	`)

	if forUpdate {
		sb.WriteString(`FOR UPDATE`)
	}

	q := sb.String()

	var consumer entity.Consumer

	err := r.dbtx.QueryRowContext(ctx, q, accountId).Scan(
		&consumer.Id,
		&consumer.AccountId,
		&consumer.IdentityNumber,
		&consumer.FullName,
		&consumer.LegalName,
		&consumer.PlaceOfBirth,
		&consumer.DateOfBirth,
		&consumer.Salary,
		&consumer.IdentityCardPhotoUrl,
		&consumer.SelfiePhotoUrl,
		&consumer.CreatedAt,
		&consumer.UpdatedAt,
		&consumer.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("[mysql_consumer_repository][GetConsumetByAccountId][QueryRowContext] error: %w | account_id: %v", err, accountId)
	}

	return &consumer, nil
}
