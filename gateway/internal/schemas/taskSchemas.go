package schemas

type CreateTaskRequest struct {
	PassData   [13][72]float64 `json:"pass_data"`
	FutureData [12][24]float64 `json:"future_data"`
	Date       string         `json:"date"`
}

type CreateTaskResponse struct {
	TaskId  int    `json:"task_id"`
	Message string `json:"message"`
}

type GetTaskRequest struct {
	TaskId int `json:"task_id" binding:"required"`
}

type GetTaskResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Date    string      `json:"date"`
	Result  [24]float64 `json:"result"`
}