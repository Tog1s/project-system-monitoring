package loadavg

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestGetFunc(t *testing.T) {
    t.Run("test get metrics", func(t *testing.T) {
        la, err := Get()

        require.NoError(t, err)
    })

}
