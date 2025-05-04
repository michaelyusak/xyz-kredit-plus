package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/apperror"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/appconstant"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/service"
)

type TransactionHandler struct {
	ctxTimeout         time.Duration
	transactionService service.TransactionService
}

func NewTransactionHandler(transctionService service.TransactionService, ctxTimeout time.Duration) *TransactionHandler {
	if ctxTimeout <= 0 {
		ctxTimeout = 30 * time.Second
	}

	return &TransactionHandler{
		ctxTimeout:         ctxTimeout,
		transactionService: transctionService,
	}
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	accountId, ok := ctx.Value(appconstant.AccountIdCtxKey).(int64)
	if !ok {
		ctx.Error(apperror.NewAppError(apperror.AppErrorOpt{
			Code:            http.StatusUnauthorized,
			ResponseMessage: http.StatusText(http.StatusUnauthorized),
		}))
		return
	}

	var req entity.Transaction

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	req.AccountId = accountId

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.ctxTimeout)
	defer cancel()

	transaction, err := h.transactionService.CreateTransaction(ctxWithTimeout, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	hHelper.ResponseOK(ctx, *transaction)
}
