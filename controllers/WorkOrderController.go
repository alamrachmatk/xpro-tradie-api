package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"log"
	"net/http"
	"strings"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

func GetAllWorkOrder(c echo.Context) error {
	var err error
	// Get parameter limit
	limitStr := c.QueryParam("limit")
	var limit uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err == nil {
			if (limit == 0) || (limit > config.LimitQuery) {
				limit = config.LimitQuery
			}
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
	} else {
		limit = config.LimitQuery
	}
	// Get parameter page
	pageStr := c.QueryParam("page")
	var page uint64
	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err == nil {
			if page == 0 {
				page = 1
			}
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
	} else {
		page = 1
	}
	var offset uint64
	if page > 1 {
		offset = limit * (page - 1)
	}
	// Get parameter pagination
	pagination := false
	paginationStr := c.QueryParam("pagination")
	if paginationStr != "" {
		pagination, err = strconv.ParseBool(paginationStr)
		if err != nil {
			return lib.CustomError(http.StatusBadRequest)
		}
	}

	var total uint64
	var responseData []models.WorkOrderList
	var httpError *echo.HTTPError
	total, responseData, httpError = WorkOrderListQuery(limit, offset, pagination)
	if httpError.Code != http.StatusOK {
		return httpError
	}

	meta := lib.GenerateMeta(c, total, limit, page, offset, pagination)
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Meta = &meta
	response.Data = responseData
	
	return c.JSON(http.StatusOK, response)

}

func WorkOrderListQuery(limit uint64, offset uint64, pagination bool) (uint64, []models.WorkOrderList, *echo.HTTPError){
	var workorders []models.WorkOrder
	total, err := models.GetAllWorkOrder(&workorders, limit, offset, pagination, nil)
	if err != nil {
		return 0, nil, lib.CustomError(http.StatusInternalServerError)
	}
	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	var responseData []models.WorkOrderList
	for _, workorder := range workorders {
		var data models.WorkOrderList
		data.WorkOrderID = workorder.WorkOrderID
		if workorder.Status == 0 {
			data.Status = "Cancel"
		} else if workorder.Status == 1 {
			data.Status = "Ongoing"
		} else if workorder.Status == 2 {
			data.Status = "Ready"
		} else if workorder.Status == 3 {
			data.Status = "Pending"
		}

		responseData = append(responseData, data)
	}

	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	return total, responseData, lib.CustomError(http.StatusOK)
}

func WorkOrderData(c echo.Context) error {
	// Get authorization token
	var token string
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	request := c.Request()
	authorization := request.Header["Authorization"]
	if authorization != nil {
		if strings.HasPrefix(authorization[0], "Bearer ") == true {
			token = authorization[0][7:]
		}
	}

	redisPool.Do("SELECT", config.RedisDBCacheToken)
	_, err := redis.String(redisPool.Do("GET", token))
	if err != nil {
		return lib.CustomError(http.StatusForbidden)
	}

	// Get parameter param
	paramStr := c.Param("id")
	log.Println("paramStr", paramStr)

	var workorder models.WorkOrder
	var data models.WorkOrderData
	status := models.GetWorkOrder(&workorder, paramStr)
	if status == 404 {
		return lib.CustomError(http.StatusNotFound)
	}

	data.WorkOrderID = workorder.WorkOrderID
	if workorder.Status == 0 {
		data.Status = "Cancel"
	} else if workorder.Status == 1 {
		data.Status = "Ongoing"
	} else if workorder.Status == 2 {
		data.Status = "Ready"
	} else if workorder.Status == 3 {
		data.Status = "Pending"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data

	return c.JSON(http.StatusOK, response)


}

func CreateWorkder(c echo.Context) error {
	// Get authorization token
	var token string
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	request := c.Request()
	authorization := request.Header["Authorization"]
	if authorization != nil {
		if strings.HasPrefix(authorization[0], "Bearer ") == true {
			token = authorization[0][7:]
		}
	}

	redisPool.Do("SELECT", config.RedisDBCacheToken)
	_, err := redis.String(redisPool.Do("GET", token))
	if err != nil {
		return lib.CustomError(http.StatusForbidden)
	}
	
	workorder := make(map[string]string)

	// Get parameter first name
	status := c.FormValue("status")
	if status != "" {
		workorder["status"] = status
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	statusResponse, _ := models.CreateWorkOrder(workorder)
	if statusResponse != 200 {
		return lib.CustomError(http.StatusInternalServerError)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)

}