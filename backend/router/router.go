package router

import (
	"eltsen00/CurrencyExchangeApp/backend/controllers"
	"eltsen00/CurrencyExchangeApp/backend/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true, // Enable cookies
		// AllowOriginFunc: func(origin string) bool { // 动态允许特定来源
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour, // 设置预检请求的缓存时间，减少预检请求的频率
	}))

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
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleByID)

		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
	}
	return r
}
