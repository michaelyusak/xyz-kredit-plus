package repository

import (
	"database/sql"
	"fmt"
)

type Transaction interface {
	Begin() error
	Rollback() error
	Commit() error
	AccountMysqlTx() *accountRepositoryMysql
	ConsumerMysqlTx() *consumerRepositoryMysql
	RefreshTokenMysqlTx() *refreshTokenRepositoryMysql
	AccountLimitMysqlTx() *accountLimitRepositoryMysql
	TransactionMysqlTx() *transactionRepositoryMysql
}

type sqlTransaction struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlTransaction(db *sql.DB) *sqlTransaction {
	return &sqlTransaction{
		db: db,
	}
}

func (s *sqlTransaction) Begin() error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("[transaction][Begin][db.Begin] Error: %w", err)
	}

	s.tx = tx

	return nil
}

func (s *sqlTransaction) Rollback() error {
	return s.tx.Rollback()
}

func (s *sqlTransaction) Commit() error {
	return s.tx.Commit()
}

func (s *sqlTransaction) AccountMysqlTx() *accountRepositoryMysql {
	return &accountRepositoryMysql{
		dbtx: s.tx,
	}
}

func (s *sqlTransaction) ConsumerMysqlTx() *consumerRepositoryMysql {
	return &consumerRepositoryMysql{
		dbtx: s.tx,
	}
}

func (s *sqlTransaction) RefreshTokenMysqlTx() *refreshTokenRepositoryMysql {
	return &refreshTokenRepositoryMysql{
		dbtx: s.tx,
	}
}

func (s *sqlTransaction) AccountLimitMysqlTx() *accountLimitRepositoryMysql {
	return &accountLimitRepositoryMysql{
		dbtx: s.tx,
	}
}

func (s *sqlTransaction) TransactionMysqlTx() *transactionRepositoryMysql {
	return &transactionRepositoryMysql{
		dbtx: s.tx,
	}
}
