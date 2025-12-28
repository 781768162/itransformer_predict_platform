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
		public.POST("/login", handler.LoginHandler)
		public.POST("/register", handler.RegisterHandler)
	}

	protected := r.Group("/api/v1/service")
	protected.Use(middleware.JwtAuth())
	{
		//JWT
		protected.POST("/create_task", handler.CreateTaskHandler)
		protected.POST("/get_task", handler.GetTaskHandler)
	}
	
	return r
}
