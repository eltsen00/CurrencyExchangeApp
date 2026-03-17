package controllers

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"eltsen00/CurrencyExchangeApp/backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateExchangeRate(c *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := c.ShouldBindJSON(&exchangeRate); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	exchangeRate.Date = time.Now()
	if err := global.Db.AutoMigrate(&models.ExchangeRate{}); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, exchangeRate)
}

func GetExchangeRates(c *gin.Context) {
	var exchangeRates []models.ExchangeRate
	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, exchangeRates)
}
