package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omarhaqqi24/marketplace-api/internal/auth"
	"github.com/omarhaqqi24/marketplace-api/internal/config"
	"github.com/omarhaqqi24/marketplace-api/internal/database"
)

func main() {
	cfg := config.Load()
	router := gin.Default()
	db := database.Connect(cfg)

	err := db.AutoMigrate(&auth.User{})
	if err != nil {
		log.Fatal(err)
	}

	repo := auth.NewUserRepository(db)
	service := auth.NewService(repo, cfg)
	handler := auth.NewHandler(service)

	router.POST("/register", handler.Register)

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
