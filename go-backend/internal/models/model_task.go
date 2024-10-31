package models

type Task struct {
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	TaskTime       string   `json:"task_time"`
	TaskFinishDate string   `json:"task_finish_date"`
	UserID         int      `json:"user_id"`
	Tags           []string `json:"tags"`
}

type User struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
