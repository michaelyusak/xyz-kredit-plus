package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/michaelyusak/go-helper/apperror"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/service"
)

type ConsumerHandler struct {
	ctxTimeout      time.Duration
	consumerService service.ConsumerService
}

func NewConsumerHandler(consumerService service.ConsumerService, ctxTimeout time.Duration) *ConsumerHandler {
	if ctxTimeout <= 0 {
		ctxTimeout = 30 * time.Second
	}

	return &ConsumerHandler{
		ctxTimeout:      ctxTimeout,
		consumerService: consumerService,
	}
}

func (h *ConsumerHandler) ProcessKyc(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

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

	identityCardPhotoFile, _, err := ctx.Request.FormFile("identity_card_photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) && consumerData.IdentityCardPhoto.Base64 == "" {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: "identity card photo is required",
			}))
			return
		}

		ctx.Error(err)
		return
	}

	selfiePhotoFile, _, err := ctx.Request.FormFile("selfie_photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) && consumerData.SelfiePhoto.Base64 == "" {
			ctx.Error(apperror.BadRequestError(apperror.AppErrorOpt{
				ResponseMessage: "selfie photo is required",
			}))
			return
		}

		ctx.Error(err)
		return
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
