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

func GetTotalBlock(c echo.Context) error {
	var responseData models.TotalBlock

	key := "totalblock"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.BlockAll = total

	return c.JSON(http.StatusOK, responseData)
}

func GetTotalDnsBlock(c echo.Context) error {
	var responseData models.TotalDnsBlock

	key := "totaldnsblock"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.DnsBlock = total

	return c.JSON(http.StatusOK, responseData)
}

func GetTotalIpAddress(c echo.Context) error {
	var responseData models.TotalIpAddress

	key := "totalipaddress"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.IpAddressAll = total

	return c.JSON(http.StatusOK, responseData)
}

func GeTotalTopMostActiveList(c echo.Context) error {
	var err error

	var totalTopMostActiveList []models.TotalTopMostActiveList
	key := "totalmostactivelist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalTopMostActiveListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalTopMostActiveListResult), &totalTopMostActiveList)
	}

	return c.JSON(http.StatusOK, totalTopMostActiveList)
}

func GeTotalDnsDayList(c echo.Context) error {
	var err error

	var totalDnsDayList []models.TotalDnsDayList
	key := "totaldnsdaylist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalDnsDayListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalDnsDayListResult), &totalDnsDayList)
	}

	return c.JSON(http.StatusOK, totalDnsDayList)
}

func GeTotalIpAddressDayList(c echo.Context) error {
	var err error

	var totalIpAddressDayList []models.TotalIpAddressDayList
	key := "totalipaddressdaylist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalDnsDayListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalDnsDayListResult), &totalIpAddressDayList)
	}

	return c.JSON(http.StatusOK, totalIpAddressDayList)
}

func GeTotalIpAddressBlockCategoryDayList(c echo.Context) error {
	var err error

	var totalIpAddressBlockCategoryDayList []models.TotalIpAddressBlockCategoryDayList
	key := "totalipaddressblockcategorydaylist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalIpAddressBlockCategoryDayListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalIpAddressBlockCategoryDayListResult), &totalIpAddressBlockCategoryDayList)
	}

	return c.JSON(http.StatusOK, totalIpAddressBlockCategoryDayList)
}

func GeTotalDnsBlockCategoryDayList(c echo.Context) error {
	var err error

	var totalDnsBlockCategoryDayList []models.TotalDnsBlockCategoryDayList
	key := "totalipaddressblockcategorydaylist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalDnsBlockCategoryDayListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalDnsBlockCategoryDayListResult), &totalDnsBlockCategoryDayList)
	}

	return c.JSON(http.StatusOK, totalDnsBlockCategoryDayList)
}
