package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureStack(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, ensureStack(nil))
	})

	t.Run("leaf-attaches-stack", func(t *testing.T) {
		leaf := New("leaf")
		got := ensureStack(leaf)
		require.NotSame(t, leaf, got)
		assert.True(t, HasStack(got))

		var stackErr StackError
		assert.True(t, As(got, &stackErr))
	})

	t.Run("already-stacked-unchanged", func(t *testing.T) {
		leaf := New("leaf")
		stacked := ensureStack(leaf)
		got := ensureStack(stacked)
		assert.Same(t, stacked, got)
		assert.True(t, HasStack(got))
	})
}
