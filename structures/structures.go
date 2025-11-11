package structures

type Task struct {
	TaskId          int    `json:"task_id"`
	TaskName        string `json:"task_name"`
	TaskDescription string `json:"task_description"`
	TaskStatus      string `json:"task_status"`
	TaskCreatedAt   string `json:"task_created_at"`
	TaskUpdatedAt   string `json:"task_updated_at"`
}
