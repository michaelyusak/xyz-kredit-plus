package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	hHelper "github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/service"
)

type TransactionHandler struct {
	ctxTimeout         time.Duration
	transactionService service.TransactionService
}

func NewTransactionHandler(ctxTimeout time.Duration, transctionService service.TransactionService) *TransactionHandler {
	if ctxTimeout <= 0 {
		ctxTimeout = 30 * time.Second
	}

	return &TransactionHandler{
		transactionService: transctionService,
	}
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var req entity.Transaction

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.ctxTimeout)
	defer cancel()

	transaction, err := h.transactionService.CreateTransaction(ctxWithTimeout, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	hHelper.ResponseOK(ctx, *transaction)
}
