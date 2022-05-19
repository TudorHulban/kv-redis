package redis

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Get handles only string values as per:
// https://redis.io/commands/get/
func (p *Pool) Get(keys ...string) ([]string, error) {
	if len(keys) == 0 {
		return nil, errors.New("no keys passed")
	}

	conn := p.pool.Get()
	defer conn.Close()

	if len(keys) == 1 {
		value, errGet := conn.Do("GET", keys[0])
		if errGet != nil {
			return nil, errGet
		}

		if value == nil {
			return nil, nil
		}

		var buf []byte
		buf = append(buf, value.([]uint8)...)

		return []string{string(buf)}, nil
	}

	// retrieve for list Redis type values,
	// TODO: move to get list function
	var res []string

	for _, key := range keys {
		value, errGet := conn.Do("LRANGE", key, 0, -1)
		if errGet != nil {
			return nil, errGet
		}

		values := value.([]interface{})

		for _, val := range values {
			res = append(res, string(val.([]uint8)))
		}
	}

	return res, nil
}

func (p *Pool) GetAny(key string, decodeInTo interface{}) error {
	conn := p.pool.Get()
	defer conn.Close()

	value, errGet := conn.Do("GET", key)
	if errGet != nil {
		return errGet
	}

	if value == nil {
		return nil
	}

	var buf []byte
	buf = append(buf, value.([]uint8)...)

	return Decoder(buf, decodeInTo)
}

// See more on SCAN in https://redis.io/commands/scan/
func (p *Pool) GetByPattern(pattern string) ([]string, error) {
	conn := p.pool.Get()
	defer conn.Close()

	var iterator int
	var keys []string

	for {
		arr, err := redis.Values(conn.Do("SCAN", iterator, "MATCH", pattern))
		if err != nil {
			return nil, fmt.Errorf("error retrieving values for '%s' pattern", pattern)
		}

		iterator, _ = redis.Int(arr[0], nil)
		iterationValues, errConv := redis.Strings(arr[1], nil)
		if errConv != nil {
			return nil, errConv
		}
		keys = append(keys, iterationValues...)

		if iterator == 0 {
			break
		}
	}

	return p.Get(keys...)
}
