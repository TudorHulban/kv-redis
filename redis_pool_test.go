package redis

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKVString(t *testing.T) {
	sock := "127.0.0.1:6379"

	pool, errNew := NewRedisPool(sock)
	require.NoError(t, errNew)
	require.NotNil(t, pool)

	now := strconv.FormatInt(time.Now().UnixNano(), 10)

	kv1 := KV{
		key:   "1" + now,
		value: "y1",
	}

	kv2 := KV{
		key:   "2" + now,
		value: "y2",
	}

	kv3 := KV{
		key:   "3" + now,
		value: "y3",
	}

	errSet := pool.Set(kv1)
	require.NoError(t, errSet, "operation set")

	pool.Set(kv2)
	pool.Set(kv3)

	value1, errGet1 := pool.Get(kv1.key)
	require.NoError(t, errGet1)
	require.Equal(t, kv1.value, value1, "get kv1")

	value2, errGet2 := pool.Get(kv2.key)
	require.NoError(t, errGet2)
	require.Equal(t, kv2.value, value2, "get kv2")

	value3, errGet3 := pool.Get(kv3.key)
	require.NoError(t, errGet3)
	require.Equal(t, kv3.value, value3, "get kv3")

	errDel1 := pool.Delete()
	require.ErrorIs(t, errNoKeysToDelete, errDel1)

	errDel2 := pool.Delete(kv1.key)
	require.NoError(t, errDel2)

	value4, errGet4 := pool.Get(kv1.key)
	require.NoError(t, errGet4)
	require.Equal(t, "", value4)

	errDel3 := pool.Delete(kv2.key, kv3.key)
	require.NoError(t, errDel3)

	value5, errGet5 := pool.Get(kv2.key)
	require.NoError(t, errGet5)
	require.Equal(t, "", value5, "kv2 should be deleted by now")

	value6, errGet6 := pool.Get(kv3.key)
	require.NoError(t, errGet6)
	require.Equal(t, "", value6, "kv3 should be deleted by now")
}

func TestKVAny(t *testing.T) {
	sock := "127.0.0.1:6379"

	pool, errNew := NewRedisPool(sock)
	require.NoError(t, errNew)
	require.NotNil(t, pool)

	v := tstruct{
		F1: 1,
		F2: []byte("a"),
	}

	key := strconv.FormatInt(time.Now().UnixNano(), 10)

	errSet := pool.SetAny(key, v)
	require.NoError(t, errSet)

	var res tstruct

	errGet := pool.GetAny(key, &res)
	require.NoError(t, errGet)
	require.Equal(t, v, res)
}
