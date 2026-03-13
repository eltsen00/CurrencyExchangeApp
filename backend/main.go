package main

import (
	"eltsen00/CurrencyExchangeApp/backend/config"
	"eltsen00/CurrencyExchangeApp/backend/router"
)

func main() {
	config.InitConfig()
	r := router.SetupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = "8080" // Default port if not specified in config
	}
	r.Run(":" + port)
}
