package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "login response",
	})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"msg": "register response",
	})
}

func UserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "userinfo response",
	})
}
