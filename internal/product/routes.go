package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler Handler,
	authMiddleware gin.HandlerFunc,
) {
	products := router.Group("/products")
	{
		products.GET("", handler.List)
		products.GET("/:id", handler.Get)

		protected := products.Group("")
		protected.Use(authMiddleware)
		protected.POST("", handler.Create)
		protected.PUT("/:id", handler.Update)
		protected.DELETE("/:id", handler.Delete)
	}
}
