package config

import (
	"encoding/json"
	"log"
	"os"
)

// General
var	Host            string = "http://localhost:1323"
var Port string = "1323"
var MediaServerPath string = "D:/images"
var MediaServerPathDl string = MediaServerPath + "/customers/driving-licence"
var MediaServerPathAvatar string = MediaServerPath + "/customers/avatar"
var	LimitQuery      uint64 = 100

// MariaDB
var	MariaDBUser     string = "root"
var	MariaDBPassword string = ""
var	MariaDBDB       string = "xpro_tradie"
var	MariaDBHost     string = "localhost"
var	MariaDBPort     string = "3306"

// Redis
var	RedisHost string = "localhost"
var	RedisPort string = "6379"

var RedisMaxIdle uint64 = 80
var RedisMaxActive uint64 = 10000
var RedisIdleTimeout uint64 = 5
var RedisWait = true

var RedisDBCacheCustomerByEmail int  = 0

type mariaDB struct {
	MariaDBUser     *string
	MariaDBPassword *string
	MariaDBDB       *string
	MariaDBHost     *string
	MariaDBPort     *string
}

type redis struct {
	RedisHost            	  	*string
	RedisPort            	  	*string
	RedisMaxIdle         		*uint64
	RedisMaxActive       		*uint64
	RedisIdleTimeout     		*uint64
	RedisWait            	   	*bool
	RedisDBCacheCustomerByEmail *int
}

type general struct {
	Host            *string
	Port            *string
	MediaServerPath *string
	LimitQuery      *uint64
}

type config struct {
	General       general
	MariaDB       mariaDB
	Redis         redis
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
	if cfg.General.MediaServerPath != nil {
		MediaServerPath = *cfg.General.MediaServerPath
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
	

	return nil
}