package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/adaptor"
	hHandler "github.com/michaelyusak/go-helper/handler"
	"github.com/michaelyusak/go-helper/helper"
	hMiddleware "github.com/michaelyusak/go-helper/middleware"
	"github.com/michaelyusak/xyz-kredit-plus/config"
	"github.com/michaelyusak/xyz-kredit-plus/handler"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
	"github.com/michaelyusak/xyz-kredit-plus/service"
	"github.com/sirupsen/logrus"
)

type routerOpts struct {
	common         *hHandler.CommonHandler
	account        *handler.AccountHandler
	allowedOrigins []string
}

func createRouter(config config.ServiceConfig, log *logrus.Logger) *gin.Engine {
	mysql, err := adaptor.ConnectDB(adaptor.MYSQL, config.MySQL)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("[server][createRouter][hAdaptor.ConnectElasticSearch] error: %s", err.Error()),
		}).Error("error connecting to ElasticSearch")
	}

	transaction := repository.NewSqlTransaction(mysql)
	accountRepo := repository.NewAccountRepositoryMysql(mysql)
	consumerRepo := repository.NewConsumerRepositoryMysql(mysql)
	RefreshTokenRepo := repository.NewRefreshTokenRepositoryMysql(mysql)

	hash := helper.NewHashHelper(config.Hash)
	jwt := helper.NewJWTHelper(config.Jwt)

	accountService := service.NewAccountService(transaction, hash, jwt, accountRepo, consumerRepo, RefreshTokenRepo)

	commonHandler := &hHandler.CommonHandler{}
	accountHandler := handler.NewAccountHandler(accountService)

	opt := routerOpts{
		common:         commonHandler,
		account:        accountHandler,
		allowedOrigins: config.AllowedOrigins,
	}

	router := newRouter(opt, log)

	return router
}

func newRouter(routerOpts routerOpts, log *logrus.Logger) *gin.Engine {
	router := gin.New()

	corsConfig := cors.DefaultConfig()

	router.ContextWithFallback = true

	router.Use(
		hMiddleware.Logger(log),
		hMiddleware.RequestIdHandlerMiddleware,
		hMiddleware.ErrorHandlerMiddleware,
		gin.Recovery(),
	)

	corsRouting(router, corsConfig, routerOpts.allowedOrigins)
	commonRouting(router, routerOpts.common)
	accountRouting(router, routerOpts.account)

	return router
}

func corsRouting(router *gin.Engine, configCors cors.Config, allowedOrigins []string) {
	configCors.AllowOrigins = allowedOrigins
	configCors.AllowMethods = []string{"POST", "GET", "PUT", "PATCH", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept", "User-Agent", "Cache-Control"}
	configCors.ExposeHeaders = []string{"Content-Length"}
	configCors.AllowCredentials = true
	router.Use(cors.New(configCors))
}

func commonRouting(router *gin.Engine, common *hHandler.CommonHandler) {
	router.GET("/ping", common.Ping)
	router.NoRoute(common.NoRoute)
}

func accountRouting(router *gin.Engine, account *handler.AccountHandler) {
	accountRouter := router.Group("/v1/account")

	accountRouter.POST("/register", account.Register)
}
