package controllers

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"eltsen00/CurrencyExchangeApp/backend/models"
	"eltsen00/CurrencyExchangeApp/backend/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashPassword

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}

	if err := global.Db.AutoMigrate(&models.User{}); err != nil {
		c.JSON(500, gin.H{"error": "Failed to migrate database"})
		return
	}

	if err := global.Db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
