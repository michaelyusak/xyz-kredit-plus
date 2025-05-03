package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/xyz-kredit-plus/appconstant"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/helper"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
)

type consumerServiceImpl struct {
	transaction  repository.Transaction
	consumerRepo repository.ConsumerRepository
	mediaRepo    repository.MediaRepository
}

func NewConsumerService(transaction repository.Transaction, consumerRepo repository.ConsumerRepository, mediaRepo repository.MediaRepository) *consumerServiceImpl {
	return &consumerServiceImpl{
		transaction:  transaction,
		consumerRepo: consumerRepo,
		mediaRepo:    mediaRepo,
	}
}

// Validate consumer data for KYC
func (s *consumerServiceImpl) validateData(consumerData entity.Consumer) error {
	// Validate data e.g. liveness, check Dukcapil, etc

	if consumerData.AccountId != 0 {
		return nil
	}

	return nil
}

// Validate consumer data
func (s *consumerServiceImpl) validateFile(media entity.Media) (repository.MediaOpt, error) {
	allowedPhotoExts := []string{".png", ".jpg"}

	var opt repository.MediaOpt

	if media.Base64 != "" {
		mediaBytes, err := base64.StdEncoding.DecodeString(media.Base64)
		if err != nil {
			return opt, fmt.Errorf("[consumer_service][validateFile][StdEncoding.DecodeString] Error: %w", err)
		}

		ext, err := helper.ValidateFileBytes(mediaBytes, allowedPhotoExts)
		if err != nil {
			return opt, fmt.Errorf("[consumer_service][validateFile][helper.ValidateFileBytes] Error: %w", err)
		}

		opt.Bytes = mediaBytes
		opt.Extension = ext

	} else if media.File != nil {
		ext, err := helper.ValidateFileMultipart(media, allowedPhotoExts)
		if err != nil {
			return opt, fmt.Errorf("[consumer_service][validateFile][helper.ValidateFileMultipart] Error: %w", err)
		}

		opt.File = &media
		opt.Extension = ext

	} else {
		return opt, fmt.Errorf("[consumer_service][validateFile] No file attached")
	}

	return opt, nil
}

func (s *consumerServiceImpl) ProcessKyc(ctx context.Context, consumerData entity.Consumer) error {
	existing, err := s.consumerRepo.GetConsumerByAccountId(ctx, consumerData.AccountId, false)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][consumerRepo.GetConsumerByAccountId] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}
	if existing != nil {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc] KYC completed | account_id: %v", consumerData.AccountId),
		})
	}

	err = s.validateData(consumerData)
	if err != nil {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			Message:         fmt.Sprintf("[consumer_service][ProcessKyc][validateData] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
			ResponseMessage: "invalid data",
		})
	}

	identityCardOpt, err := s.validateFile(consumerData.IdentityCardPhoto)
	if err != nil {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			Message:         fmt.Sprintf("[consumer_service][ProcessKyc][validateFile][identityCard] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
			ResponseMessage: "invalid or corrupted indentity card photo",
		})
	}

	consumerData.IdentityCardPhoto.Key = helper.HashSHA256(fmt.Sprintf("%v%s", consumerData.AccountId, appconstant.KYCIdentityCardPhotoTag))

	selfiePhotoOpt, err := s.validateFile(consumerData.SelfiePhoto)
	if err != nil {
		return apperror.BadRequestError(apperror.AppErrorOpt{
			Message:         fmt.Sprintf("[consumer_service][ProcessKyc][validateFile][selfie] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
			ResponseMessage: "invalid or corrupted selfie photo",
		})
	}

	consumerData.SelfiePhoto.Key = helper.HashSHA256(fmt.Sprintf("%v%s", consumerData.AccountId, appconstant.KYCSelfiePhotoTag))

	err = s.transaction.Begin()
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][transaction.Begin] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}

	consumerRepo := s.transaction.ConsumerMysqlTx()

	defer func() {
		if err != nil {
			s.transaction.Rollback()
		}

		s.transaction.Commit()
	}()

	_, err = consumerRepo.GetConsumerByAccountId(ctx, consumerData.AccountId, true)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][consumerRepo.GetConsumerByAccountId][ForUpdate] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}

	err = consumerRepo.InsertConsumer(ctx, consumerData)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][consumerRepo.InsertConsumer] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}

	identityCardOpt.Key = consumerData.IdentityCardPhoto.Key
	err = s.mediaRepo.Store(ctx, identityCardOpt)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][mediaRepo.Store][identityCard] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}

	selfiePhotoOpt.Key = consumerData.SelfiePhoto.Key
	err = s.mediaRepo.Store(ctx, selfiePhotoOpt)
	if err != nil {
		return apperror.InternalServerError(apperror.AppErrorOpt{
			Message: fmt.Sprintf("[consumer_service][ProcessKyc][mediaRepo.Store][selfie] Error: %s | account_id: %v", err.Error(), consumerData.AccountId),
		})
	}

	return nil
}
