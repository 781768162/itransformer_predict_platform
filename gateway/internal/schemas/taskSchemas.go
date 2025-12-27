package schemas

type CreateTaskRequest struct {
	PassData [13][72]float64
	FutureData [12][24]float64
}

type CreateTaskResponse struct {
	TaskId int `json:"task_id"`
	Message string `json:"message"`
}

type GetTaskRequest struct {
	TaskId int `json:"task_id" binding:"required"`
}

type GetTaskResponse struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Result [24]float64 `json:"result"`
}