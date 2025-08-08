package gql

type ProjectData struct {
	Project Project `env:"project"`
}

type Project struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Environments Environments `json:"environments"`
	Services     Services     `json:"services"`
}

type Environments struct {
	Edges []struct {
		Node struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"node"`
	} `json:"edges"`
}

type Services struct {
	Edges []ServiceEdge `json:"edges"`
}

type ServiceEdge struct {
	Node Service `json:"node"`
}

type Service struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	ServiceInstances ServiceInstances `json:"serviceInstances"`
}

type ServiceInstances struct {
	Edges []ServiceInstanceEdge `json:"edges"`
}

type ServiceInstanceEdge struct {
	Node ServiceInstance `json:"node"`
}

type ServiceInstance struct {
	Id               string     `json:"id"`
	ServiceId        string     `json:"serviceId"`
	EnvironmentId    string     `json:"environmentId"`
	LatestDeployment Deployment `json:"latestDeployment"`
}

type Deployment struct {
	CanRedeploy bool           `json:"canRedeploy"`
	Id          string         `json:"id"`
	Meta        DeploymentMeta `json:"meta"`
}

type DeploymentMeta struct {
	ServiceManifest ServiceManifest `json:"serviceManifest"`
}

type ServiceManifest struct {
	Deploy DeployConfig `json:"deploy"`
}

type DeployConfig struct {
	MultiRegionConfig map[string]RegionConfig `json:"multiRegionConfig"`
}

type RegionConfig struct {
	NumReplicas int `json:"numReplicas"`
}
