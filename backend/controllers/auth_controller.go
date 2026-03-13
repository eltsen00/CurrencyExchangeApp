package controllers

import (
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
	c.JSON(200, gin.H{"token": token})
}
