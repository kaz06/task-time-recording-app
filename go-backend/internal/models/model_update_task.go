package models

type UpdateTask struct {
	Title          string   `json:"title"`
	TaskTime       string   `json:"task_time"`
	TaskFinishDate string   `json:"task_finish_date"`
	UserID         int      `json:"user_id"`
	Tags           []string `json:"tags"`
}
