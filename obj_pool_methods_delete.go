package redis

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

func (p *RedisPool) deleteByDB() error {
	conn := p.pool.Get()
	defer conn.Close()

	_, errDel := conn.Do("FLUSHDB", "ASYNC")
	return errDel
}

func (p *RedisPool) deleteALL() error {
	conn := p.pool.Get()
	defer conn.Close()

	_, errDel := conn.Do("FLUSHALL", "ASYNC")
	return errDel
}
