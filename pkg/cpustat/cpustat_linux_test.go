package cpustat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFunc(t *testing.T) {
	t.Run("test get cpustat", func(t *testing.T) {
		cs, err := Get()

		require.NoError(t, err)
		require.NotNil(t, cs.User)
		require.NotNil(t, cs.System)
		require.NotNil(t, cs.Idle)
	})
}
