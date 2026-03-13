package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Login successful"})
			})
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Registration successful"})
			})
		}
	}
	return r
}
