package api

import "github.com/ferretcode/switchyard/incident/pkg/types"

type ListIncidentReportsResponse struct {
	DeploymentReports []types.IncidentReport `json:"deployment_reports"`
	GeneralReports    []types.IncidentReport `json:"general_reports"`
}
