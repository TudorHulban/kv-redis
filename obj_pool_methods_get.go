package redis

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

func (p *RedisPool) GetAny(key string, decodeInTo interface{}) error {
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
