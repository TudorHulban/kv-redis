package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

type KV struct {
	key   string
	value string
}

type RedisPool struct {
	pool redis.Pool
}

var errNoKeysToDelete = errors.New("no keys to delete")

func NewRedisPool(sock string) (*RedisPool, error) {
	var errConn error

	res := RedisPool{
		pool: redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", sock)
				if err != nil {
					errConn = err
				}

				return c, err
			},
		},
	}

	if errConn != nil {
		return nil, errConn
	}

	return &res, nil
}

// See more about SET in:
// https://redis.io/commands/set/
func (p *RedisPool) Set(kv KV) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SET", kv.key, kv.value)
	return errSet
}

// Get handles only string values as per:
// https://redis.io/commands/get/
func (p *RedisPool) Get(key string) (string, error) {
	conn := p.pool.Get()
	defer conn.Close()

	value, errGet := conn.Do("GET", key)
	if errGet != nil {
		return "", errGet
	}

	if value == nil {
		return "", nil
	}

	var buf []byte
	buf = append(buf, value.([]uint8)...)

	return string(buf), nil
}

func (p *RedisPool) Delete(keys ...string) error {
	conn := p.pool.Get()
	defer conn.Close()

	if len(keys) == 0 {
		return errNoKeysToDelete
	}

	if len(keys) == 1 {
		_, errDel := conn.Do("DEL", keys[0])
		return errDel
	}

	redisKeys := make([]interface{}, len(keys))

	for i := 0; i < len(keys); i++ {
		redisKeys[i] = interface{}(keys[i])
	}

	_, errDel := conn.Do("DEL", redisKeys...)
	return errDel
}
