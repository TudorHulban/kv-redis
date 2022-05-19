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

type Pool struct {
	pool                redis.Pool
	maxNumberNamespaces uint
	databaseNumber      uint
}

type PoolOption func(p *Pool)

var errNoKeysToDelete = errors.New("no keys to delete")

func WithDatabaseNumber(n uint) PoolOption {
	return func(p *Pool) {
		if n > p.maxNumberNamespaces {
			p.databaseNumber = p.maxNumberNamespaces
		} else {
			p.databaseNumber = n
		}
	}
}

func NewPool(sock string, config ...PoolOption) (*Pool, error) {
	res := Pool{
		maxNumberNamespaces: 16,
	}

	for _, option := range config {
		option(&res)
	}

	var errConn error

	res.pool = redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", sock)
			if err != nil {
				errConn = err
			}

			return c, err
		},
	}

	if errConn != nil {
		return nil, errConn
	}

	return &res, nil
}

func (p *Pool) Close() {
	p.pool.Close()
}
