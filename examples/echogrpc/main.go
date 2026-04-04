package main

import (
	"github.com/jkaveri/goservice"
	"github.com/jkaveri/goservice/examples/echogrpc/service"
)

func main() {
	goservice.Run(service.New)
}
