package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("wrong file in path", func(t *testing.T) {
		_, err := ReadFromFile("./configs/config.toml")
		require.Error(t, err)
	})
}
