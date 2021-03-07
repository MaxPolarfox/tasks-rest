package types

type Task struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type TaskIdResponse struct {
	ID string `json:"id"`
}

type AddTaskReqBody struct {
	Data string `json:"data"`
}
