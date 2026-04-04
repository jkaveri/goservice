package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/env"
)

func Test_IsLocal(t *testing.T) {
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, rEnv.IsDebug, env.IsLocal())
}

func Test_IsProduction(t *testing.T) {
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, rEnv.DeploymentEnv == env.PROD, env.IsProduction())
}

func Test_IsDev(t *testing.T) {
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, rEnv.DeploymentEnv == env.DEV, env.IsDev())
}

func Test_IsQA(t *testing.T) {
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, rEnv.DeploymentEnv == env.QA, env.IsQA())
}

func Test_SetDeploymentEnv(t *testing.T) {
	env.SetDeploymentEnv(env.DEV)
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, env.DEV, rEnv.DeploymentEnv)
}
