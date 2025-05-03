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
			identity_card_photo_key, 
			selfie_photo_key, 
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
		&consumer.IdentityCardPhoto.Key,
		&consumer.SelfiePhoto.Key,
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

func (r *consumerRepositoryMysql) InsertConsumer(ctx context.Context, consumerData entity.Consumer) error {
	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO consumers (account_id, identity_number, full_name, legal_name, place_of_birth, date_of_birth, salary, identity_card_photo_key, selfie_photo_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	q := sb.String()

	now := nowUnixMilli()

	_, err := r.dbtx.ExecContext(ctx, q,
		consumerData.AccountId,
		consumerData.IdentityNumber,
		consumerData.FullName,
		consumerData.LegalName,
		consumerData.PlaceOfBirth,
		consumerData.DateOfBirth,
		consumerData.Salary,
		consumerData.IdentityCardPhoto.Key,
		consumerData.SelfiePhoto.Key,
		now,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}
