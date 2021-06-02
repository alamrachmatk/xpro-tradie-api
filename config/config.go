package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// General
var Host string = "http://localhost:1323"
var Port string = "1323"
var LimitQuery uint64 = 100
var ReadTimeout time.Duration = 5
var WriteTimeout time.Duration = 10
var IdleTimeout time.Duration = 120

// MariaDB
var MariaDBUser string = "root"
var MariaDBPassword string = "quantum"
var MariaDBDB string = "quantumdns"
var MariaDBHost string = "175.106.13.14"
var MariaDBPort string = "3306"

// Redis
var RedisHost string = "localhost"
var RedisPort string = "6379"

var RedisMaxIdle uint64 = 80
var RedisMaxActive uint64 = 10000
var RedisIdleTimeout uint64 = 5
var RedisWait = true

var RedisDBCacheToken int = 0
var RedisDBCacheCustomerByEmail int = 1

// Cache
var ExpireToken uint64 = 3600 // 1 hour

// JWT xspr0Tr@Die
var JWTSecret string = "SL6ANV4cMfu2cBI240iV0xYLgv6RxUIh"

type mariaDB struct {
	MariaDBUser     *string
	MariaDBPassword *string
	MariaDBDB       *string
	MariaDBHost     *string
	MariaDBPort     *string
}

type redis struct {
	RedisHost                   *string
	RedisPort                   *string
	RedisMaxIdle                *uint64
	RedisMaxActive              *uint64
	RedisIdleTimeout            *uint64
	RedisWait                   *bool
	RedisDBCacheCustomerByEmail *int
}

type general struct {
	Host                  *string
	Port                  *string
	MediaServerPath       *string
	MediaServerPathDl     *string
	MediaServerPathAvatar *string
	LimitQuery            *uint64
	// ReadTimeout covers the time from when the connection is accepted to
	// when the request body is fully read
	ReadTimeout *time.Duration
	// WriteTimeout normally covers the time from the end of the request
	// header read to the end of the response write
	WriteTimeout *time.Duration
	// IdleTimeout which limits server-side the amount of time a Keep-Alive
	// connection will be kept idle before being reused
	IdleTimeout *time.Duration
}

type jwt struct {
	JWTSecret *string
}

type config struct {
	General general
	MariaDB mariaDB
	Redis   redis
	JWT     jwt
}

func InitConfig() error {
	// configPtr := flag.String("c", "../config/metubeapi.json", "meTube Core API configuration file")

	file, err := os.Open("../config/xpro_tradie.json")
	if err != nil {
		log.Println("Use default configuration")
		return nil
	}
	decoder := json.NewDecoder(file)
	var cfg config
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	// General
	if cfg.General.Host != nil {
		Host = *cfg.General.Host
	}
	if cfg.General.Port != nil {
		Port = *cfg.General.Port
	}
	if cfg.General.ReadTimeout != nil {
		ReadTimeout = *cfg.General.ReadTimeout
	}
	if cfg.General.WriteTimeout != nil {
		WriteTimeout = *cfg.General.WriteTimeout
	}
	if cfg.General.IdleTimeout != nil {
		IdleTimeout = *cfg.General.IdleTimeout
	}

	// MariaDB
	if cfg.MariaDB.MariaDBDB != nil {
		MariaDBDB = *cfg.MariaDB.MariaDBDB
	}
	if cfg.MariaDB.MariaDBHost != nil {
		MariaDBHost = *cfg.MariaDB.MariaDBHost
	}
	if cfg.MariaDB.MariaDBPort != nil {
		MariaDBPort = *cfg.MariaDB.MariaDBPort
	}

	// Redis
	if cfg.Redis.RedisHost != nil {
		RedisHost = *cfg.Redis.RedisHost
	}
	if cfg.Redis.RedisPort != nil {
		RedisPort = *cfg.Redis.RedisPort
	}
	if cfg.Redis.RedisMaxIdle != nil {
		RedisMaxIdle = *cfg.Redis.RedisMaxIdle
	}
	if cfg.Redis.RedisMaxActive != nil {
		RedisMaxActive = *cfg.Redis.RedisMaxActive
	}
	if cfg.Redis.RedisIdleTimeout != nil {
		RedisIdleTimeout = *cfg.Redis.RedisIdleTimeout
	}
	if cfg.Redis.RedisWait != nil {
		RedisWait = *cfg.Redis.RedisWait
	}
	if cfg.Redis.RedisDBCacheCustomerByEmail != nil {
		RedisDBCacheCustomerByEmail = *cfg.Redis.RedisDBCacheCustomerByEmail
	}

	// JWT
	if cfg.JWT.JWTSecret != nil {
		JWTSecret = *cfg.JWT.JWTSecret
	}

	return nil
}
