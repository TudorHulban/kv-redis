package redis

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const _sock = "127.0.0.1:6379"

func TestKVString(t *testing.T) {
	pool, errNew := NewPool(_sock)
	require.NoError(t, errNew)
	require.NotNil(t, pool)

	pool.deleteByDB()
	defer pool.Close()

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

	require.NoError(t, pool.Set(kv1), "operation set kv1")
	require.NoError(t, pool.Set(kv2), "operation set kv2")
	require.NoError(t, pool.Set(kv3), "operation set kv3")

	value1, errGet1 := pool.Get(kv1.key)
	require.NoError(t, errGet1)
	require.Equal(t, []string{kv1.value}, value1, "get kv1")

	value2, errGet2 := pool.Get(kv2.key)
	require.NoError(t, errGet2)
	require.Equal(t, []string{kv2.value}, value2, "get kv2")

	value3, errGet3 := pool.Get(kv3.key)
	require.NoError(t, errGet3)
	require.Equal(t, []string{kv3.value}, value3, "get kv3")

	require.ErrorIs(t, errNoKeysToDelete, pool.Delete())

	require.NoError(t, pool.Delete(kv1.key))

	value4, errGet4 := pool.Get(kv1.key)
	require.NoError(t, errGet4)
	require.Nil(t, value4)

	require.NoError(t, pool.Delete(kv2.key, kv3.key))

	value5, errGet5 := pool.Get(kv2.key)
	require.NoError(t, errGet5)
	require.Nil(t, value5, "kv2 should be deleted by now")

	value6, errGet6 := pool.Get(kv3.key)
	require.NoError(t, errGet6)
	require.Nil(t, value6, "kv3 should be deleted by now")
}

func TestKVAny(t *testing.T) {
	pool, errNew := NewPool(_sock)
	require.NoError(t, errNew)
	require.NotNil(t, pool)

	pool.deleteByDB()
	defer pool.Close()

	v := tstruct{
		F1: 1,
		F2: []byte("a"),
	}

	key := strconv.FormatInt(time.Now().UnixNano(), 10)

	require.NoError(t, pool.SetAny(key, v))

	var res1 tstruct

	require.NoError(t, pool.GetAny(key, &res1))
	require.Equal(t, v, res1)

	var res2 tstruct

	require.NoError(t, pool.GetAny("1", &res2))
	require.Equal(t, tstruct{}, res2, "no key therefore no change in passed var")
}

func TestKVs(t *testing.T) {
	pool, errNew := NewPool(_sock)
	require.NoError(t, errNew)
	require.NotNil(t, pool)

	pool.deleteByDB()
	defer pool.Close()

	kv1 := KVs{
		key:    "K1",
		values: []string{"1", "2"},
	}

	kv2 := KVs{
		key:    "K2",
		values: []string{"3"},
	}

	require.NoError(t, pool.SetList(kv1))
	require.NoError(t, pool.SetList(kv2))

	values, errGet := pool.GetByPattern("*")
	require.NoError(t, errGet)
	require.Equal(t, 3, len(values), fmt.Sprintf("values should have 2 items, instead: %s", values))
}
