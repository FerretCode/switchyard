package scheduler

type RegisterWorkerServiceRequest struct {
	ServiceId string `json:"service_id"`
	JobName   string `json:"job_name"`
}

type ScheduleJobRequest struct {
	JobName    string         `json:"job_name"`
	JobContext map[string]any `json:"job_context"`
}
