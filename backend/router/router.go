package router

import (
	"eltsen00/CurrencyExchangeApp/backend/controllers"
	"eltsen00/CurrencyExchangeApp/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/exchangeRates", controllers.GetExchangeRates)
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}
	}
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
	}
	return r
}
