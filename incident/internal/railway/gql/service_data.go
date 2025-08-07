package gql

type ServiceData struct {
	Service Service `json:"service"`
}

type Service struct {
	Id          string      `json:"id"`
	Deployments Deployments `json:"deployments"`
}

type Deployments struct {
	Edges []Edge `json:"edges"`
}

type Edge struct {
	Node Node `json:"node"`
}

type Node struct {
	Id              string `json:"id"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	StatusUpdatedAt string `json:"statusUpdatedAt"`
	Status          string `json:"status"`
}
