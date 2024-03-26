package loadavg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("test get function", func(t *testing.T) {
		la, err := Get()

		require.NoError(t, err)
		require.NotNil(t, la.LoadAvg1)
		require.NotNil(t, la.LoadAvg5)
		require.NotNil(t, la.LoadAvg15)
	})
}
