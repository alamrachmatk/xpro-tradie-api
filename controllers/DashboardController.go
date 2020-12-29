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

func GetTotalSite(c echo.Context) error {
	var responseData models.TotalSite

	key := "totalsites"
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	total, err := redis.Uint64(redisPool.Do("GET", key+".value"))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.SiteAll = total

	return c.JSON(http.StatusOK, responseData)
}

func GeTotalTopMostActiveList(c echo.Context) error {
	var err error

	var totalTopMostActiveList []models.TotalTopMostActiveList
	// responseData, _ = TotalTopMostActiveListQuery(limit)
	key := "mostactive"
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

// func TotalTopMostActiveListQuery(limit uint64) ([]models.TotalTopMostActiveList, *echo.HTTPError) {
// 	var mostactives []models.MostActive
// 	err := models.GetTotalTopMostActiveListQuery(&mostactives, limit)
// 	if err != nil {
// 		return nil, lib.CustomError(http.StatusInternalServerError)
// 	}

// 	var responseData []models.TotalTopMostActiveList
// 	for _, mostactive := range mostactives {
// 		var data models.TotalTopMostActiveList
// 		data.BaseDomain = *mostactive.BaseDomain
// 		data.Total = *mostactive.TotalMostActive

// 		responseData = append(responseData, data)
// 	}

// 	return responseData, nil
// }
