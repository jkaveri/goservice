package goservice

import "fmt"

func tcpAddrFromPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
