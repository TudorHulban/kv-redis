package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

type KV struct {
	key   string
	value string
}

type KVs struct {
	key    string
	values []string
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

func (p *RedisPool) Close() {
	p.pool.Close()
}
