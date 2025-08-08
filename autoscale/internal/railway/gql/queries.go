package gql

import _ "embed"

//go:embed project.graphql
var ProjectQuery string

//go:embed update_regions.graphql
var UpdateRegionsQuery string

//go:embed metrics.graphql
var MetricsQuery string

//go:embed service_instance_deploy.graphql
var ServiceInstanceDeployQuery string
