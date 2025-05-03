package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/michaelyusak/go-helper/apperror"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/helper"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
)

type accountServiceImpl struct {
	transaction      repository.Transaction
	hash             hHelper.HashHelper
	jwt              hHelper.JWTHelper
	accountRepo      repository.AccountRepository
	consumerRepo     repository.ConsumerRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAccountService(transaction repository.Transaction, hash hHelper.HashHelper, jwt hHelper.JWTHelper, accountRepo repository.AccountRepository, consumerRepo repository.ConsumerRepository, refreshTokenRepo repository.RefreshTokenRepository) *accountServiceImpl {
	return &accountServiceImpl{
		transaction:      transaction,
		hash:             hash,
		jwt:              jwt,
		accountRepo:      accountRepo,
		consumerRepo:     consumerRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *accountServiceImpl) generateJwt(account entity.Account, isKycCompleted bool) (*entity.TokenData, error) {
	customClaims := make(map[string]any)
	customClaims["account_id"] = account.Id
	customClaims["email"] = account.Email
	customClaims["is_kyc_completed"] = isKycCompleted

	claimsBytes, err := json.Marshal(customClaims)
	if err != nil {
		return nil, fmt.Errorf("[account_service][generateJwt][json.Marshal] Error: %w", err)
	}

	accessTokenExpiredAt := time.Now().Add(time.Hour).UnixMilli()

	accessToken, err := s.jwt.CreateAndSign(claimsBytes, accessTokenExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("[account_service][generateJwt][jwt.CreateAndSign][accessToken] Error: %w", err)
	}

	refreshTokenExpiredAt := time.Now().Add(24 * time.Hour).UnixMilli()

	refreshToken, err := s.jwt.CreateAndSign(claimsBytes, refreshTokenExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("[account_service][generateJwt][jwt.CreateAndSign][refreshToken] Error: %w", err)
	}

	return &entity.TokenData{
		AccessToken: entity.Token{
			Token:     accessToken,
			ExpiredAt: accessTokenExpiredAt,
		},
		RefreshToken: entity.Token{
			Token:     refreshToken,
			ExpiredAt: refreshTokenExpiredAt,
		},
	}, nil
}

func (s *accountServiceImpl) RegisterAccount(ctx context.Context, newAccount entity.Account) (*entity.TokenData, error) {
	if !helper.ValidatePassword(newAccount.Password) {
		return nil, apperror.BadRequestError(apperror.AppErrorOpt{
			ResponseMessage: "invalid password",
		})
	}

	err := s.transaction.Begin()
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][transaction.Begin] Error: %s", err.Error()),
		})
	}

	accountRepo := s.transaction.AccountMysqlTx()
	consumerRepo := s.transaction.ConsumerMysqlTx()
	refreshTokenRepo := s.transaction.RefreshTokenMysqlTx()

	defer func() {
		if err != nil {
			s.transaction.Rollback()
		}

		s.transaction.Commit()
	}()

	existing, err := accountRepo.GetAccountByEmail(ctx, newAccount.Email, true)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][accountRepo.GetAccountByEmail] Error: %s", err.Error()),
		})
	}
	if existing != nil {
		return nil, apperror.BadRequestError(apperror.AppErrorOpt{
			Message:         "[account_service][Register] email already registered",
			ResponseMessage: "email already registered",
		})
	}

	hashed, err := s.hash.Hash(newAccount.Password)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][hash.Hash] Error: %s", err.Error()),
		})
	}

	newAccount.Password = hashed

	accountId, err := accountRepo.InsertAccount(ctx, newAccount)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][accountRepo.InsertAccount] Error: %s", err.Error()),
		})
	}

	newAccount.Id = accountId

	existingConsumer, err := consumerRepo.GetConsumerByAccountId(ctx, newAccount.Id, false)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][consumerRepo.GetConsumetByAccountId] Error: %s", err.Error()),
		})
	}

	isKycCompleted := false

	if existingConsumer != nil {
		isKycCompleted = true
	}

	token, err := s.generateJwt(newAccount, isKycCompleted)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][generateJwt] Error: %s", err.Error()),
		})
	}

	err = refreshTokenRepo.InsertToken(ctx, token.RefreshToken.Token, newAccount.Id, token.RefreshToken.ExpiredAt)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][RegisterAccount][refreshTokenRepo.InsertToken] Error: %s", err.Error()),
		})
	}

	return token, nil
}

func (s *accountServiceImpl) Login(ctx context.Context, account entity.Account) (*entity.TokenData, error) {
	existing, err := s.accountRepo.GetAccountByEmail(ctx, account.Email, false)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][Login][accountRepo.GetAccountByEmail] Error: %s | email: %s", err.Error(), account.Email),
		})
	}
	if existing == nil {
		return nil, apperror.NewAppError(apperror.AppErrorOpt{
			Code:            http.StatusUnauthorized,
			Message:         "[account_service][Login] email not registered",
			ResponseMessage: "email not registered",
		})
	}

	isValid, err := s.hash.Check(account.Password, []byte(existing.Password))
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][Login][hash.Check] Error: %s | email: %s", err.Error(), account.Email),
		})
	}
	if !isValid {
		return nil, apperror.NewAppError(apperror.AppErrorOpt{
			Code:            http.StatusUnauthorized,
			Message:         "[account_service][Login] invalid credentials",
			ResponseMessage: "invalid credentials",
		})
	}

	existingConsumer, err := s.consumerRepo.GetConsumerByAccountId(ctx, existing.Id, false)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][Login][consumerRepo.GetConsumerByAccountId] Error: %s | account_id: %v", err.Error(), existing.Id),
		})
	}

	isKycCompleted := false

	if existingConsumer != nil {
		isKycCompleted = true
	}

	token, err := s.generateJwt(account, isKycCompleted)
	if err != nil {
		return nil, apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[account_service][Login][generateJwt] Error: %s | account_id: %v", err.Error(), existing.Id),
		})
	}

	return token, nil
}
