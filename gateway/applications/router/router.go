package router

import (
	"gateway/internal/handler"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func MustNewRouter() *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1/user")
	{
		// 无需鉴权
		public.POST("/user/login", handler.LoginHandler)
		public.POST("/user/register", handler.RegisterHandler)
	}

	protected := r.Group("/api/v1/service")
	protected.Use(middleware.JwtAuth())
	{
		//JWT
	}
	
	return r
}
