package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler Handler,
	authMiidleware gin.HandlerFunc,
) {
	products := router.Group("/products")
	{
		products.GET("", handler.List)
		products.GET("/:id", handler.Get)

		protected := products.Group("")
		protected.Use(authMiidleware)
		protected.POST("", handler.Create)
	}
}
