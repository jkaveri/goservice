module github.com/jkaveri/goservice/examples/echogrpc

go 1.25.0

replace (
	github.com/jkaveri/goservice => ../../
	github.com/jkaveri/goservice/logger/zap => ../../logger/zap
)

require (
	github.com/jkaveri/goservice v0.0.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	git.toolsfdg.net/fe/golog v0.0.0-20260305064040-9dc0a7d30d68 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jkaveri/goconfig v1.0.0 // indirect
	github.com/jkaveri/ramda v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
)
