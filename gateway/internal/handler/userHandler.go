package handler

import (
	"net/http"

	"gateway/internal/code"
	"gateway/internal/logic"
	"gateway/internal/schemas"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var req schemas.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("ShouldBindJSON error: %v", err)
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	token, expireAt, err := logic.LoginLogic(c.Request.Context(), req.UserName, req.Password)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.LoginResponse{
		Token:      token,
		ExpireTime: expireAt,
		Message:    "success",
	}
	c.JSON(http.StatusOK, resp)
}

func RegisterHandler(c *gin.Context) {
	var req schemas.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("ShouldBindJSON error: %v", err)
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	err := logic.RegisterLogic(c.Request.Context(), req.UserName, req.Password)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.RegisterResponse{
		Message: "created",
	}
	c.JSON(http.StatusOK, resp)
}
