package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/env"
)

func Test_IsDebug(t *testing.T) {
	rEnv := env.GetRuntimeInfo()

	assert.Equal(t, rEnv.IsDebug, env.IsDebug())
}
