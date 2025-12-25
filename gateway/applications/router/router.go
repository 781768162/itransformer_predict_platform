package router

import (
	"gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

func MustNewRouter() *gin.Engine {
	r := gin.Default()

	router := r.Group("/api/v1")
	{
		router.POST("/user/login", handler.Login)
		router.POST("/user/register", handler.Register)
		router.GET("/user/userinfo", handler.UserInfo)

	}
	return r
}
