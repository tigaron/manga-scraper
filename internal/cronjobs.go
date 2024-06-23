package internal

type CronJob struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Crontab string `json:"crontab"`
	Tags    string `json:"tags"`
}

type CronJobStatus struct {
	ID       string  `json:"id"`
	JobID    string  `json:"jobId"`
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Duration float64 `json:"duration"`
}

type CreateCronJobParams struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Crontab string `json:"crontab"`
	Tags    string `json:"tags"`
}

type CreateCronJobStatusParams struct {
	JobID    string  `json:"jobId"`
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Duration float64 `json:"duration"`
}

type UpdateCronJobStatusParams struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Duration float64 `json:"duration"`
}
