package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/darzox/test_task_iq_progress/docs"
)

type Handler interface {
	Deposit(ctx *gin.Context)
	Transfer(ctx *gin.Context)
	GetLast10Transactions(ctx *gin.Context)
}

func RegisterRoutes(handler Handler, router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	router.POST("/deposit", handler.Deposit)
	router.POST("/transfer", handler.Transfer)
	router.GET("/transactions", handler.GetLast10Transactions)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
