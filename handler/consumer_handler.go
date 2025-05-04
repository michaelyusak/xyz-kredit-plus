package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/michaelyusak/go-helper/apperror"
	_ "github.com/michaelyusak/go-helper/dto"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/appconstant"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/service"
)

type ConsumerHandler struct {
	ctxTimeout      time.Duration
	consumerService service.ConsumerService
}

func NewConsumerHandler(consumerService service.ConsumerService, ctxTimeout time.Duration) *ConsumerHandler {
	if ctxTimeout <= 0 {
		ctxTimeout = 2 * time.Minute
	}

	return &ConsumerHandler{
		ctxTimeout:      ctxTimeout,
		consumerService: consumerService,
	}
}

// Consumer godoc
// @Summary Process a KYC for an account
// @Description Post consumer data for KYC, including personal information and photos (identity card and selfie).
// @Description Consumer JSON structure: see model entity.Consumer
// @Description Example Data: {"nik": "124","full_name": "user test","legal_name": "user test legal","place_of_birth": "bumi","date_of_birth": "12-07-2001","salary": 600000,"identity_card_photo": {"base64":"image_base64_encoded"},"selfie_photo": {"base64": "image_base64_encoded"}}
// @Tags consumers
// @Accept  multipart/form-data
// @Produce  json
// @Param Authorization header string true "Bearer token"
// @Param data formData string true "Consumer data in JSON format (same as entity.Consumer)" example={"nik": "124","full_name": "user test","legal_name": "user test legal","place_of_birth": "bumi","date_of_birth": "12-07-2001","salary": 600000,"identity_card_photo": {"base64":""},"selfie_photo": {"base64": ""}}
// @Param identity_card_photo formData file false "Identity card photo"
// @Param selfie_photo formData file false "Selfie photo"
// @Success 200 {object} dto.Response{message=string,data=nil} "Success"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or validation error"
// @Router /consumer/process-kyc [post]
func (h *ConsumerHandler) ProcessKyc(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, ok := ctx.Value(appconstant.AccountIdCtxKey).(int64)
	if !ok {
		ctx.Error(apperror.NewAppError(apperror.AppErrorOpt{
			Code:            http.StatusUnauthorized,
			ResponseMessage: http.StatusText(http.StatusUnauthorized),
		}))
		return
	}

	var consumerData entity.Consumer

	data := ctx.Request.FormValue("data")
	if data == "" {
		ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
			ResponseMessage: "data cannot be empty",
		}))
		return
	}

	err := json.Unmarshal([]byte(data), &consumerData)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = validator.New().Struct(consumerData)
	if err != nil {
		ctx.Error(err)
		return
	}

	consumerData.AccountId = accountId

	fileSizeLimit := 3 * 1024 * 1024 // 3MB

	identityCardPhotoFile, identityCardPhotoHeader, err := ctx.Request.FormFile("identity_card_photo")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			ctx.Error(err)
			return
		}

		if consumerData.IdentityCardPhoto.Base64 == "" {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: "identity card photo is required",
			}))
			return
		}
	}
	if identityCardPhotoFile != nil && identityCardPhotoHeader != nil {
		if identityCardPhotoHeader.Size > int64(fileSizeLimit) {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: fmt.Sprintf(
					"identity card photo file size %.2fMB exceeded limit of %s",
					float32(identityCardPhotoHeader.Size/(1024*1024)),
					"3MB"),
			}))
		}
	}

	selfiePhotoFile, selfiePhotoHeader, err := ctx.Request.FormFile("selfie_photo")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			ctx.Error(err)
			return
		}

		if consumerData.SelfiePhoto.Base64 == "" {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: "selfie photo is required",
			}))
			return
		}
	}
	if selfiePhotoFile != nil && selfiePhotoHeader != nil {
		if selfiePhotoHeader.Size > int64(fileSizeLimit) {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: fmt.Sprintf(
					"selfie photo file size %.2fMB exceeded limit of %s",
					float32(selfiePhotoHeader.Size/(1024*1024)),
					"3MB"),
			}))
		}
	}

	consumerData.IdentityCardPhoto.File = identityCardPhotoFile
	consumerData.SelfiePhoto.File = selfiePhotoFile

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.ctxTimeout)
	defer cancel()

	err = h.consumerService.ProcessKyc(ctxWithTimeout, consumerData)
	if err != nil {
		ctx.Error(err)
		return
	}

	hHelper.ResponseOK(ctx, nil)
}
