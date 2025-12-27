package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gateway/internal/code"
	"gateway/internal/logic"
	"gateway/internal/schemas"
	"gateway/pkg/parse"

	"github.com/gin-gonic/gin"
)

func CreateTaskHandler(c *gin.Context) {
	f, err := c.FormFile("file")
	date := c.PostForm("date")
	if err != nil || len(date) == 0 {
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	os.MkdirAll("/tmp/uploads", 0755)
	dst := filepath.Join("/tmp/uploads", fmt.Sprintf("%d_%s", time.Now().UnixNano(), f.Filename))
	if err := c.SaveUploadedFile(f, dst); err != nil {
		status, body := code.ToHTTP(code.ErrSaveFile)
		c.JSON(status, body)
		return
	}

	passData, futureData, err := parse.ParseCSV(dst, date) // 解析csv
	if err != nil {
		status, body := code.ToHTTP(code.ErrParseFile)
		c.JSON(status, body)
		return
	}

	userIdVal, ok := c.Get("userId")
	if !ok {
		status, body := code.ToHTTP(code.ErrUnknown)
		c.JSON(status, body)
		return
	}
	userId, _ := userIdVal.(int64)
	
	taskId, err := logic.CreateTaskLogic(c.Request.Context(), userId, passData, futureData)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.CreateTaskResponse{
		TaskId: taskId,
		Message: "success",
	}
	c.JSON(http.StatusOK, resp)
}

func GetTaskHandler(c *gin.Context) {
	var req schemas.GetTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	status, res, err := logic.GetTaskLogic(c.Request.Context(), req.TaskId)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.GetTaskResponse{
		Message: "",
		Status: status,
		Result: res,
	}
	c.JSON(http.StatusOK, resp)
}

