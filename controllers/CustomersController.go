package controllers

import (
	"api/config"
	"api/db"
	"encoding/json"
	"api/lib"
	"api/models"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

func SignUp(c echo.Context) error  {
	params := make(map[string]string)

	// Get parameter first name
	firstName := c.FormValue("first_name")
	if firstName != "" {
		params["first_name"] = firstName
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter last name
	lastName := c.FormValue("last_name")
	if lastName != "" {
		params["last_name"] = lastName
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter email
	email := c.FormValue("email")
	if email != "" {
		params["email"] = email
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter phone
	phone := c.FormValue("phone")
	if phone != "" {
		params["phone"] = phone
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter address
	address := c.FormValue("address")
	if address != "" {
		params["address"] = address
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter category
	category := c.FormValue("category")
	if category != "" {
		params["category"] = category
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter company name
	companyName := c.FormValue("company_name")
	if companyName != "" {
		params["company_name"] = companyName
	}

	// Get parameter abn/cn number
	abnCnNumber := c.FormValue("abn_cn_number")
	if abnCnNumber != "" {
		params["abn_cn_number"] = abnCnNumber
	}

	// Get parameter driving licence
	drivingLicence := c.FormValue("driving_licence")
	if drivingLicence != "" {
		params["driving_licence"] = drivingLicence
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter photo id
	photoId := c.FormValue("photo_id")
	if photoId != "" {
		params["photo_id"] = photoId
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter avatar
	avatar := c.FormValue("avatar")
	if avatar != "" {
		params["avatar"] = avatar
	}

	// Get parameter avatar
	status := c.FormValue("status")
	if status != "" {
		params["status"] = status
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	redisPool.Do("SELECT", config.RedisDBCacheCustomerByEmail)
	_, err := redis.Bytes(redisPool.Do("GET", email))
	if err != nil {
		statusResponse := models.CreateCustomer(params)
		if statusResponse != 200 {
			return lib.CustomError(http.StatusInternalServerError)
		}
		dataJSON, _ := json.Marshal(params)
		redisPool.Do("SET", email, dataJSON)
	} else {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Email already in use")
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)

}
