package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gateway/config"
	"gateway/internal/code"
	"gateway/internal/logic"
	"gateway/internal/schemas"
	"gateway/pkg/parse"
	"gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CreateTaskHandler(c *gin.Context) {
	date := c.PostForm("date")
	if date == "" {
		logger.Errorf("PostForm date is empty")
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	passFile, err := c.FormFile("passdata_file")
	if err != nil {
		logger.Errorf("FormFile err: %v", err)
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	futureFile, err := c.FormFile("futuredata_file")
	if err != nil {
		logger.Errorf("FormFile err: %v", err)
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	err = os.MkdirAll(config.Settings.Server.UploadDir, 0755)
	if err != nil {
		logger.Errorf("MkdirAll err: %v", err)
		status, body := code.ToHTTP(code.ErrParseFile)
		c.JSON(status, body)
		return
	}

	passDst := filepath.Join(config.Settings.Server.UploadDir, fmt.Sprintf("%d_pass_%s", time.Now().UnixNano(), passFile.Filename))
	futureDst := filepath.Join(config.Settings.Server.UploadDir, fmt.Sprintf("%d_future_%s", time.Now().UnixNano(), futureFile.Filename))

	if err := c.SaveUploadedFile(passFile, passDst); err != nil {
		logger.Errorf("SaveUploadedFile passDst: %s err: %v", passDst, err)
		status, body := code.ToHTTP(code.ErrSaveFile)
		c.JSON(status, body)
		return
	}
	if err := c.SaveUploadedFile(futureFile, futureDst); err != nil {
		logger.Errorf("SaveUploadedFile futureDst: %s err: %v", futureDst, err)
		status, body := code.ToHTTP(code.ErrSaveFile)
		c.JSON(status, body)
		return
	}

	passData, futureData, err := parse.ParseCSV(passDst, futureDst)
	if err != nil {
		logger.Errorf("ParseCSV err: %v", err)
		status, body := code.ToHTTP(code.ErrParseFile)
		c.JSON(status, body)
		return
	}

	userIdVal, ok := c.Get("userId")
	if !ok {
		logger.Errorf("Get userId is empty")
		status, body := code.ToHTTP(code.ErrUnknown)
		c.JSON(status, body)
		return
	}
	userId, ok := userIdVal.(int64)
	if !ok {
		logger.Errorf("userIdVal is Invalid")
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	taskId, err := logic.CreateTaskLogic(c.Request.Context(), userId, date, passData, futureData)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.CreateTaskResponse{
		TaskId:  taskId,
		Message: "success",
	}
	c.JSON(http.StatusOK, resp)
}

func GetTaskHandler(c *gin.Context) {
	var req schemas.GetTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("ShouldBindJSON error: %v", err)
		status, body := code.ToHTTP(code.ErrInvalidParam)
		c.JSON(status, body)
		return
	}

	status, date, res, err := logic.GetTaskLogic(c.Request.Context(), req.TaskId)
	if err != nil {
		status, body := code.ToHTTP(err)
		c.JSON(status, body)
		return
	}

	resp := &schemas.GetTaskResponse{
		Message: "",
		Status:  status,
		Date:    date,
		Result:  res,
	}
	c.JSON(http.StatusOK, resp)
}
