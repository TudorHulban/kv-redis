package redis

import "fmt"

// See more about SET in:
// https://redis.io/commands/set/
func (p *RedisPool) Set(kv KV) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SET", kv.key, kv.value)
	return errSet
}
func (p *RedisPool) SetAny(key string, any interface{}) error {
	buf, errEnc := Encoder(any)
	if errEnc != nil {
		return fmt.Errorf("set any: %w", errEnc)
	}

	conn := p.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SET", key, string(buf))
	return errSet
}
