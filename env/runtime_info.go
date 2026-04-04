package env

type RuntimeInfo struct {
	IsDebug       bool
	DeploymentEnv string
	ServiceName   string
	Version       string
}

func GetRuntimeInfo() RuntimeInfo {
	return runtimeInfo
}

func GetDeploymentEnv() string {
	return runtimeInfo.DeploymentEnv
}

func SetVersion(s string) {
	runtimeInfo.Version = s
}

func GetVersion() string {
	return runtimeInfo.Version
}

func SetServiceName(s string) {
	runtimeInfo.ServiceName = s
}

func GetServiceName() string {
	return runtimeInfo.ServiceName
}
