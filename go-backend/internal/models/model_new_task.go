package models

type NewTask struct {
	Title          string   `json:"title"`
	TaskTime       string   `json:"task_time"`
	TaskFinishDate string   `json:"task_finish_date"`
	Tags           []string `json:"tags"`
}
