package messagebus

const (
	OK = iota
	ERROR
)

type FinishJobMessage struct {
	JobId   string `json:"job_id"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
