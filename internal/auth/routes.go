package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler Handler,
) {
	auths := router.Group("/auth")
	{
		auths.POST("/register", handler.Register)
		auths.POST("/login", handler.Login)
	}
}
