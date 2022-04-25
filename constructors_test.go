package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	sock := "127.0.0.1:6379"

	c, close, errNew := NewRedisClient(sock, 1)
	require.NoError(t, errNew)

	require.NotNil(t, c)
	close()
}
