package env

func IsLocal() bool {
	return runtimeInfo.DeploymentEnv == LOCAL
}

func IsDev() bool {
	return runtimeInfo.DeploymentEnv == DEV
}

func IsQA() bool {
	return runtimeInfo.DeploymentEnv == QA
}

func IsProduction() bool {
	return runtimeInfo.DeploymentEnv == PROD
}

func SetDeploymentEnv(env string) {
	runtimeInfo.DeploymentEnv = env
}
