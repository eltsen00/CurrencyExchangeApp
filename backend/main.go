package main

import (
	"context"
	"eltsen00/CurrencyExchangeApp/backend/config"
	"eltsen00/CurrencyExchangeApp/backend/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.InitConfig()
	r := router.SetupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = "8080" // Default port if not specified in config
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // 监听CMD中断信号
	<-quit                            // 阻塞，直到接收到中断信号
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置5秒的超时时间，确保服务器在关闭前完成当前请求
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	} else {
		log.Println("Server exited properly")
	}
}
