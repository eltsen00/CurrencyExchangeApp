package controllers

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"eltsen00/CurrencyExchangeApp/backend/models"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	AllcacheKey   = "articles_cache"
	cacheDuration = 10 * time.Minute
	IDcacheKey    = "article_id_cache"
)

func CreateArticle(c *gin.Context) {
	var article models.Article

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&models.Article{}); err != nil {
		c.JSON(500, gin.H{"error": "Failed to migrate database"})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create article"})
		return
	}
	if err := global.RedisDB.Del(AllcacheKey).Err(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear article cache"})
		return
	}

	c.JSON(201, article)
}

func GetArticles(c *gin.Context) {

	cacheData, err := global.RedisDB.Get(AllcacheKey).Result()

	if errors.Is(err, redis.Nil) { // Cache miss

		var articles []models.Article

		// Fetch articles from the database
		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"error": "No articles found"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to retrieve articles"})
			}
			return
		}

		// Cache the articles in Redis
		jsonData, err := json.Marshal(articles)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to serialize articles"})
			return
		}
		err = global.RedisDB.Set(AllcacheKey, jsonData, cacheDuration).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to cache articles"})
			return
		}

		c.JSON(200, articles)
		return
	} else if err != nil { // Redis error
		c.JSON(500, gin.H{"error": "Failed to retrieve articles from cache"})
		return
	} else { // Cache hit
		var articles []models.Article
		err = json.Unmarshal([]byte(cacheData), &articles)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to deserialize articles"})
			return
		}
		c.JSON(200, articles)
		return
	}
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	cacheKey := IDcacheKey + id
	cacheData, err := global.RedisDB.Get(cacheKey).Result()

	if errors.Is(err, redis.Nil) { // Cache miss
		var article models.Article

		// Fetch article from the database
		if err := global.Db.First(&article, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"error": "Article not found"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to retrieve article"})
			}
			return
		}
		// Cache the article in Redis
		jsonData, err := json.Marshal(article)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to serialize article"})
			return
		}
		err = global.RedisDB.Set(cacheKey, jsonData, cacheDuration).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to cache article"})
			return
		}

		c.JSON(200, article)
		return
	} else if err != nil { // Redis error
		c.JSON(500, gin.H{"error": "Failed to retrieve article from cache"})
		return
	} else { // Cache hit
		var article models.Article
		err = json.Unmarshal([]byte(cacheData), &article)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to deserialize article"})
			return
		}
		c.JSON(200, article)
		return
	}
}
