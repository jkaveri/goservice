module github.com/jkaveri/goservice/examples/grpcgateway

go 1.25.0

replace (
	github.com/jkaveri/goservice => ../../
	github.com/jkaveri/goservice/logger/zap => ../../logger/zap
)

require (
	git.toolsfdg.net/fe/golog v1.1.0
	github.com/jkaveri/goservice v0.0.0
	github.com/envoyproxy/protoc-gen-validate v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/genproto/googleapis/api v0.0.0-20260316180232-0b37fe3546d5
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jkaveri/goconfig v1.0.0 // indirect
	github.com/jkaveri/ramda v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260316180232-0b37fe3546d5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
