module github.com/gigmcp/registry

go 1.26.2

require github.com/gigmcp/registry/schema v0.0.0

require (
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace github.com/gigmcp/registry/schema => ./schema
