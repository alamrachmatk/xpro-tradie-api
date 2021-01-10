package controllers

import (
	"api/db"
	"api/lib"
	"api/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

func GetTotalDns(c echo.Context) error {
	var responseData models.TotalDns

	key := "totaldns"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.DnsAll = total

	return c.JSON(http.StatusOK, responseData)
}

func GetTotalBlok(c echo.Context) error {
	var responseData models.TotalBlok

	key := "totalblok"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.BlokAll = total

	return c.JSON(http.StatusOK, responseData)
}

func GetTotalDnsBlok(c echo.Context) error {
	var responseData models.TotalDnsBlok

	key := "totaldnsblok"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.DnsBlok = total

	return c.JSON(http.StatusOK, responseData)
}

func GeTotalTopMostActiveList(c echo.Context) error {
	var err error

	var totalTopMostActiveList []models.TotalTopMostActiveList
	// responseData, _ = TotalTopMostActiveListQuery(limit)
	key := "mostactivelist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalTopMostActiveResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalTopMostActiveResult), &totalTopMostActiveList)
	}

	return c.JSON(http.StatusOK, totalTopMostActiveList)
}
