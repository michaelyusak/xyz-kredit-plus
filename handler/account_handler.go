package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/service"
)

type AccountHandler struct {
	ctxTimeout     time.Duration
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Register(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var newAccount entity.Account

	err := ctx.ShouldBind(&newAccount)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.ctxTimeout)
	defer cancel()

	token, err := h.accountService.RegisterAccount(ctxWithTimeout, newAccount)
	if err != nil {
		ctx.Error(err)
		return
	}

	helper.ResponseOK(ctx, *token)
}
