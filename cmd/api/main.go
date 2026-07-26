package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/omarhaqqi24/marketplace-api/internal/config"
	"github.com/omarhaqqi24/marketplace-api/internal/database"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)
	router := gin.Default()

	_ = db

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":    cfg.AppName,
			"env":    cfg.AppEnv,
			"status": "ok",
		})
	})

	fmt.Println("App is running on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
