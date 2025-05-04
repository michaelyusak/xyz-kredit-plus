package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type transactionRepositoryMysql struct {
	dbtx DBTX
}

func NewTransactionRepositoryMysql(dbtx DBTX) *transactionRepositoryMysql {
	return &transactionRepositoryMysql{
		dbtx: dbtx,
	}
}

func (r *transactionRepositoryMysql) InsertTransaction(ctx context.Context, transaction entity.Transaction) (int64, error) {
	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO transactions (
    		account_id,
    		contact_number,
    		otr,
    		admin_fee,
    		total_installment,
    		total_interest,
    		asset_name,
    		created_at,
    		updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	q := sb.String()

	now := nowUnixMilli()

	res, err := r.dbtx.ExecContext(ctx, q,
		transaction.AccountId,
		transaction.ContactNumber,
		transaction.OTR,
		transaction.AdminFee,
		transaction.TotalInstallemnt,
		transaction.TotalInterest,
		transaction.AssetName,
		now,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("[mysql_transaction_repository][InsertTransaction][ExecContext] error: %w | account_id: %v", err, transaction.AccountId)
	}

	transactionId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[mysql_transaction_repository][InsertTransaction][LastInsertId] error: %w | account_id: %v", err, transaction.AccountId)
	}


	return transactionId, nil
}
