package autoscale

type RegisterServiceRequest struct {
	ServiceId string `json:"service_id"`
	JobName   string `json:"job_name"`
}
