package autoscale

type RegisterServiceRequest struct {
	ServiceId string `json:"service_id"`
	JobName   string `json:"job_name"`
}

type ListServicesResponse struct {
	Services []ServiceContext `json:"services"`
}

type ServiceContext struct {
	ServiceId              string  `json:"service_id"`
	ProjectId              string  `json:"project_id"`
	JobName                string  `json:"job_name"`
	EnvironmentName        string  `json:"environment_name"`
	EnvironmentId string  `json:"environment_id"`
	Replicas               int     `json:"replicas"`
	MinReplicas            int     `json:"min_replicas"`
	MaxReplicas            int     `json:"max_replicas"`
	CpuUpscaleThreshold    float64 `json:"cpu_upscale_threshold"`
	MemoryUpscaleThreshold float64 `json:"memory_upscale_threshold"`
	UpscaleCooldown        string  `json:"upscale_cooldown"`
	DownscaleCooldown      string  `json:"downscale_cooldown"`
	LastScaledAt           string  `json:"last_scaled_at"`
}
