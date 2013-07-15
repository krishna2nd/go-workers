package workers

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type config struct {
	processId string
	pool      *redis.Pool
}

var Config *config

func Configure(options map[string]string) {
	var poolSize int

	if options["server"] == "" {
		panic("Configure requires a 'server' option, which identifies a Redis instance")
	}
	if options["process"] == "" {
		panic("Configure requires a 'process' option, which uniquely identifies this instance")
	}
	if options["pool"] == "" {
		options["pool"] = "1"
	}

	poolSize, _ = strconv.Atoi(options["pool"])

	Config = &config{
		options["process"],
		&redis.Pool{
			MaxIdle:     poolSize,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", options["server"])
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}
