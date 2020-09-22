package db

import (
	"api/config"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var Redis redis.Conn
var RedisPool *redis.Pool

func RedisInit() redis.Conn {
	conn, err := redis.Dial("tcp", config.RedisHost+":"+config.RedisPort)
	if err != nil {
		log.Fatal("Error connect to Redis server!")
	}
	return conn
}

// From https://stackoverflow.com/questions/24387350/re-using-redigo-connection-instead-of-recreating-it-every-time
func RedisPoolInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     int(config.RedisMaxIdle),
		MaxActive:   int(config.RedisMaxActive), // max number of connections
		IdleTimeout: time.Duration(config.RedisIdleTimeout),
		Wait:        config.RedisWait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisHost+":"+config.RedisPort)
			if err != nil {
				log.Fatal("Error connect to Redis server!")
			}
			return c, err
		},
	}
}
