package controllers

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(c *gin.Context) {
	articleID := c.Param("id")

	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to like article"})
		return
	}
	c.JSON(200, gin.H{"message": "Article liked successfully"})
}

func GetArticleLikes(c *gin.Context) {
	articleID := c.Param("id")
	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			likes = "0"
		} else {
			c.JSON(500, gin.H{"error": "Failed to get article likes"})
			return
		}
	}
	c.JSON(200, gin.H{"likes": likes})
}
