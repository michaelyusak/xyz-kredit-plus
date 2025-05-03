package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/michaelyusak/go-helper/adaptor"
	hHandler "github.com/michaelyusak/go-helper/handler"
	hMiddleware "github.com/michaelyusak/go-helper/middleware"
	"github.com/michaelyusak/xyz-kredit-plus/config"
	"github.com/sirupsen/logrus"
)

type routerOpts struct {
	common         *hHandler.CommonHandler
	allowedOrigins []string
}

func createRouter(config config.ServiceConfig, log *logrus.Logger) *gin.Engine {
	_, err := adaptor.ConnectDB(adaptor.MYSQL, config.MySQL)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("[server][createRouter][hAdaptor.ConnectElasticSearch] error: %s", err.Error()),
		}).Error("error connecting to ElasticSearch")
	}

	commonHandler := &hHandler.CommonHandler{}

	opt := routerOpts{
		common:         commonHandler,
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
