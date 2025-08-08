package gql

import _ "embed"

//go:embed variable_collection_upsert.graphql
var VariableCollectionUpsertQuery string

//go:embed project.graphql
var ProjectQuery string
