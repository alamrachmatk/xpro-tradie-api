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

func GeTotalRequestList(c echo.Context) error {
	var err error

	var totalRequestList []models.TotalRequestList
	key := "totalrequestlist"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	totalRequestListResult, err := redis.Bytes(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusNotFound)
	} else {
		json.Unmarshal([]byte(totalRequestListResult), &totalRequestList)
	}

	return c.JSON(http.StatusOK, totalRequestList)
}
