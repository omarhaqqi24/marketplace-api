package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omarhaqqi24/marketplace-api/internal/auth"
	"github.com/omarhaqqi24/marketplace-api/internal/config"
	"github.com/omarhaqqi24/marketplace-api/internal/database"
	"github.com/omarhaqqi24/marketplace-api/internal/product"
)

func main() {
	cfg := config.Load()
	router := gin.Default()
	db := database.Connect(cfg)

	err := db.AutoMigrate(
		&auth.User{},
		&product.Product{},
	)

	if err != nil {
		log.Fatal(err)
	}

	authRepo := auth.NewUserRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	productRepo := product.NewProductRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":    cfg.AppName,
			"env":    cfg.AppEnv,
			"status": "ok",
		})
	})

	auth.RegisterRoutes(router, authHandler)
	product.RegisterRoutes(router, productHandler, auth.AuthMiddleware(cfg.JWTSecret))

	fmt.Println("App is running on port", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
