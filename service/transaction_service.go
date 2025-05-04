package service

import (
	"context"
	"fmt"

	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
)

type transactionServiceImpl struct {
	transaction      repository.Transaction
	accountLimitRepo repository.AccountLimitRepository
	transactionRepo  repository.TransactionRepository
}

func NewTransactionService(transaction repository.Transaction, accountLimitRepos repository.AccountLimitRepository, transactionRepo repository.TransactionRepository) *transactionServiceImpl {
	return &transactionServiceImpl{
		transaction:      transaction,
		accountLimitRepo: accountLimitRepos,
		transactionRepo:  transactionRepo,
	}
}

func (s *transactionServiceImpl) adjustLimit(limit entity.AccountLimit, otr float64, installemntMonths int) entity.AccountLimit {
	var lim float64

	switch installemntMonths {
	case 1:
		lim = limit.Limit1M

	case 2:
		lim = limit.Limit2M

	case 3:
		lim = limit.Limit3M

	case 4:
		lim = limit.Limit4M
	}

	discount := (lim - otr) / lim

	newLimit := limit

	newLimit.Limit1M = limit.Limit1M * discount
	newLimit.Limit2M = limit.Limit2M * discount
	newLimit.Limit3M = limit.Limit3M * discount
	newLimit.Limit4M = limit.Limit4M * discount

	return newLimit
}

func (s *transactionServiceImpl) CreateTransaction(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error) {
	err := s.transaction.Begin()
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[transaction_service][CreateTransaction][transaction.Begin] Error: %s | account_id: %v", err.Error(), transaction.AccountId),
		})
	}

	accountLimitRepo := s.transaction.AccountLimitMysqlTx()
	transactionRepo := s.transaction.TransactionMysqlTx()

	defer func() {
		if err != nil {
			s.transaction.Rollback()
		}

		s.transaction.Commit()
	}()

	limit, err := accountLimitRepo.GetAccountLimitByAccountId(ctx, transaction.AccountId, true)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[transaction_service][CreateTransaction][accountLimitRepo.GetAccountLimitByAccountId] Error: %s | account_id: %v", err.Error(), transaction.AccountId),
		})
	}

	isLimitSufficient := false

	switch transaction.InstallmentMonths {
	case 1:
		if transaction.OTR <= limit.Limit1M {
			isLimitSufficient = true
		}

	case 2:
		if transaction.OTR <= limit.Limit2M {
			isLimitSufficient = true
		}

	case 3:
		if transaction.OTR <= limit.Limit3M {
			isLimitSufficient = true
		}

	case 4:
		if transaction.OTR <= limit.Limit4M {
			isLimitSufficient = true
		}
	default:
		return nil, apperror.BadRequestError(apperror.AppErrorOpt{
			Message: "maximum installemnt months is 4",
		})
	}

	if !isLimitSufficient {
		return nil, apperror.BadRequestError(apperror.AppErrorOpt{
			ResponseMessage: "insufficient limit",
		})
	}

	newLimit := s.adjustLimit(*limit, transaction.OTR, transaction.InstallmentMonths)

	err = accountLimitRepo.UpdateLimit(ctx, newLimit)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[transaction_service][CreateTransaction][accountLimitRepo.UpdateLimit] Error: %s | account_id: %v", err.Error(), transaction.AccountId),
		})
	}

	transactionId, err := transactionRepo.InsertTransaction(ctx, transaction)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[transaction_service][CreateTransaction][transactionRepo.InsertTransaction] Error: %s | account_id: %v", err.Error(), transaction.AccountId),
		})
	}

	transaction.Id = transactionId

	return &transaction, nil
}
