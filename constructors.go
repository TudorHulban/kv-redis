package redis

import (
	"fmt"

	"github.com/rueian/rueidis"
)

func NewRedisClient(sock string, namespace uint) (rueidis.Client, func(), error) {
	if namespace > 16 {
		return nil, nil, fmt.Errorf("provided namespace; %d is out of range\n", namespace)
	}

	c, _ := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{sock},
		SelectDB:    int(namespace),
	})

	closeClient := func() {
		c.Close()
	}

	return c, closeClient, nil
}
