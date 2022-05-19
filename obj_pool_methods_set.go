package redis

import "fmt"

// See more about SET in:
// https://redis.io/commands/set/
func (p *Pool) Set(kv KV) error {
	conn := p.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SET", kv.key, kv.value)
	return errSet
}

func (p *Pool) SetAny(key string, any interface{}) error {
	buf, errEnc := Encoder(any)
	if errEnc != nil {
		return fmt.Errorf("set any: %w", errEnc)
	}

	conn := p.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SET", key, string(buf))
	return errSet
}

func (p *Pool) SetList(kv KVs) error {
	conn := p.pool.Get()
	defer conn.Close()

	var serializedItems []interface{}
	serializedItems = append(serializedItems, kv.key)

	for _, item := range kv.values {
		serializedItems = append(serializedItems, item)
	}

	_, errSet := conn.Do("RPUSH", serializedItems...)
	return errSet
}
