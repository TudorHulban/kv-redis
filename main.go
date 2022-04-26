package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool = newPool()

func main() {
	client := pool.Get()
	defer client.Close()

	_, err := client.Do("SET", "mykey", "Hi from redigo!")
	if err != nil {
		panic(err)
	}

	value, err := client.Do("GET", "mykey")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s \n", value)

	_, err = client.Do("ZADD", "vehicles", 4, "car")
	if err != nil {
		panic(err)
	}
	_, err = client.Do("ZADD", "vehicles", 2, "bike")
	if err != nil {
		panic(err)
	}

	vehicles, err := client.Do("ZRANGE", "vehicles", 0, -1, "WITHSCORES")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s \n", vehicles)

}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
