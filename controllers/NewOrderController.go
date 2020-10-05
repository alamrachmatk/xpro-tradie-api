package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"net/http"
	"strings"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

func CreateNewOrder(c echo.Context) error {
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

	params := make(map[string]string)

	// Get parameter name
	name := c.FormValue("name")
	if name != "" {
		params["name"] = name
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter customer id
	customerId := c.FormValue("customer_id")
	if customerId != "" {
		params["customer_id"] = customerId
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter company setting ID
	companySettingID := c.FormValue("company_setting_id")
	if companySettingID != "" {
		params["company_setting_id"] = companySettingID
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter due date
	dueDate := c.FormValue("due_date")
	if dueDate != "" {
		params["due_date"] = dueDate
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter budget
	budget := c.FormValue("budget")
	if budget != "" {
		params["budget"] = budget
	}

	// Get parameter description
	description := c.FormValue("description")
	if description != "" {
		params["description"] = description
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter status
	status := c.FormValue("status")
	if status != "" {
		params["status"] = status
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	timeLayout := "2006-01-02 15:04:05"
	timeNow := time.Now()
	pareseDueDate, _ := time.Parse(timeLayout, dueDate)
	if timeNow.Format("2006-01-02 15:04:05") >= pareseDueDate.Format("2006-01-02 15:04:05") {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "The due date is less than the current date")
	}
	
	statusResponse, _ := models.CreateNewOrder(params)
	if statusResponse != 200 {
		return lib.CustomError(http.StatusInternalServerError)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)
}