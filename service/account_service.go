package service

import (
	"context"
	"fmt"

	"github.com/michaelyusak/go-helper/apperror"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/helper"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
)

type accountServiceImpl struct {
	transaction repository.Transaction
	hash        hHelper.HashHelper
	accountRepo repository.AccountRepository
}

func NewAccountService(transaction repository.Transaction, hash hHelper.HashHelper, accountRepo repository.AccountRepository) *accountServiceImpl {
	return &accountServiceImpl{
		transaction: transaction,
		hash:        hash,
		accountRepo: accountRepo,
	}
}

func (s *accountServiceImpl) RegisterAccount(ctx context.Context, newAccount entity.Account) error {
	if !helper.ValidatePassword(newAccount.Password) {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			ResponseMessage: "invalid password",
		})
	}

	err := s.transaction.Begin()
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][transaction.Begin] Error: %s", err.Error()),
		})
	}

	accountRepo := s.transaction.AccountMysqlTx()

	defer func() {
		if err != nil {
			s.transaction.Rollback()
		}

		s.transaction.Commit()
	}()

	err = accountRepo.Lock(ctx)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][accountRepo.Lock] Error: %s", err.Error()),
		})
	}

	existing, err := accountRepo.GetAccountByEmail(ctx, newAccount.Email)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][accountRepo.GetAccountByEmail] Error: %s", err.Error()),
		})
	}
	if existing != nil {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			Message:         "[account_service][Register] email already registered",
			ResponseMessage: "email already registered",
		})
	}

	hashed, err := s.hash.Hash(newAccount.Password)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][hash.Hash] Error: %s", err.Error()),
		})
	}

	newAccount.Password = hashed

	err = accountRepo.InsertAccount(ctx, newAccount)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][accountRepo.InsertAccount] Error: %s", err.Error()),
		})
	}

	return nil
}
