package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/env"
)

func Test_HostNameToServiceName(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		sn := env.HostNameToServiceName("dc-identity-v2-6984f87b96-2mn7q")
		assert.Equal(t, "dc-identity-v2", sn)
	})
}
