package router

import (
	"os"
	"path/filepath"

	"gateway/internal/handler"
	"gateway/internal/middleware"
	

	"github.com/gin-gonic/gin"
)

func MustNewRouter() *gin.Engine {
	r := gin.Default()

	// 静态资源目录基于可执行文件所在路径
	exeDir, err := os.Executable()
	staticDir := "./static"
	if err == nil {
		staticDir = filepath.Join(filepath.Dir(exeDir), "static")
	}

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

	r.Static("/", staticDir)

	return r
}
