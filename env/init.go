package env

import (
	"os"
	"strings"

	"github.com/jkaveri/ramda"
)

var (
	ServiceName string
	Version     string
	runtimeInfo RuntimeInfo
)

func init() {
	runtimeInfo = RuntimeInfo{
		IsDebug:       strings.EqualFold(os.Getenv("DEBUG"), "true"),
		DeploymentEnv: ramda.Default(os.Getenv("DEPLOYMENT_ENV"), DEV),
		ServiceName: ramda.DefaultFn(func() string {
			if sv, exist := os.LookupEnv("SERVICE_NAME"); exist && sv != "" {
				return sv
			}

			hn, _ := os.Hostname()

			return HostNameToServiceName(hn)
		}, ServiceName),
		Version: ramda.Default(Version, "unknown"),
	}
}
