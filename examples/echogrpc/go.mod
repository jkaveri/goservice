module github.com/jkaveri/goservice/examples/echogrpc

go 1.26.1

replace github.com/jkaveri/goservice => ../../

require (
	github.com/jkaveri/goservice v0.0.0
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jkaveri/goconfig v1.0.0 // indirect
	github.com/jkaveri/golog/v2 v2.1.0 // indirect
	github.com/jkaveri/ramda v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260316180232-0b37fe3546d5 // indirect
)
