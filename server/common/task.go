package common

// 任务相关

type Task struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	Bash string `json:"bash"`
}

