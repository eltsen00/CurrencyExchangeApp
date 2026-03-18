package controllers

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"eltsen00/CurrencyExchangeApp/backend/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	c.JSON(201, article)
}

func GetArticles(c *gin.Context) {
	var articles []models.Article

	if err := global.Db.Find(&articles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "No articles found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve articles"})
		}
		return
	}

	c.JSON(200, articles)
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := global.Db.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Article not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve article"})
		}
		return
	}
	c.JSON(200, article)
}
