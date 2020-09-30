package controllers

import (
	"api/config"
	"api/db"
	"api/models"
	"api/lib"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"log"
	 
	"github.com/labstack/echo"
	"github.com/garyburd/redigo/redis"
)

func CustomerData(c echo.Context) error  {
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

	tokenMap, err := lib.ExtractClaims(token)
	if err != nil {
		return lib.CustomError(http.StatusBadGateway)
	} 
  
	var customer models.CustomerDataCache
	email := tokenMap["email"]
	redisPool.Do("SELECT", config.RedisDBCacheCustomerByEmail)
	dataToken, err := redis.Bytes(redisPool.Do("GET", email))
	if err != nil { 
		return lib.CustomError(http.StatusForbidden)
	}
	err = json.Unmarshal(dataToken, &customer)
	if err != nil {
		log.Println("error unmarshalProperties:", err)
		return lib.CustomError(http.StatusInternalServerError,
			"Sorry, we've experience a problem. Please try again.",
			"Internal server error")
	}

	var data models.CustomerData
	CustomerID, _ := strconv.ParseUint(customer.CustomerID, 10, 64)
	data.CustomerID = CustomerID
	data.FirstName = customer.FirstName
	data.LastName = customer.LastName
	data.Email = customer.Email
	data.Phone = customer.Phone
	data.Address = customer.Address
	if customer.Category == "1" {
		data.Category = "Company"
	} else {
		data.Category = "Housing"
	}
	data.CompanyName = customer.CompanyName
	data.AbnCnNumber = customer.AbnCnNumber
	data.DrivingLicence = customer.DrivingLicence
	data.PhotoId = customer.PhotoId
	if customer.Avatar != nil {
		data.Avatar = customer.Avatar
	} else {

	}
	if customer.Status == "0" {
		data.Status = "New"
	} else if customer.Status == "1" {
		data.Status = "Active"
	} else if customer.Status == "2" {
		data.Status = "Banned"
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Data = data
	
	return c.JSON(http.StatusOK, response)
}

func UpdateCustomerData(c echo.Context) error {
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

	tokenMap, err := lib.ExtractClaims(token)
	if err != nil {
		return lib.CustomError(http.StatusBadGateway)
	} 
  
	idCache := tokenMap["customer_id"]

	params := make(map[string]string)
	idStr := c.Param("id")
	if idStr != "" {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if id > 0 {
			params["id"] = strconv.FormatUint(id,10)
			var customer models.Customer
			status := models.GetCustomer(&customer, idStr)
			if status != http.StatusOK {
				return lib.CustomError(http.StatusNotFound)
			}
		}
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	if idCache != idStr {
		return lib.CustomError(http.StatusForbidden)
	}

	// Get parameter first name
	firstName := c.FormValue("first_name")
	if firstName != "" {
		params["first_name"] = firstName
	}

	// Get parameter last name
	lastName := c.FormValue("last_name")
	if lastName != "" {
		params["last_name"] = lastName
	}

	// Get parameter phone
	phone := c.FormValue("phone")
	if phone != "" {
		params["phone"] = phone
	}

	// Get parameter address
	address := c.FormValue("address")
	if address != "" {
		params["address"] = address
	}

	// Get parameter category
	category := c.FormValue("category")
	if category != "" {
		params["category"] = category
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

	// Get parameter photo id
	photoId := c.FormValue("photo_id")
	if photoId != "" {
		params["photo_id"] = photoId
	}

	// Get parameter avatar
	status := c.FormValue("status")
	if status != "" {
		params["status"] = status
	}

	var fileDl *multipart.FileHeader
	fileDl, err = c.FormFile("driving_licence")
	if fileDl != nil {
		// Get file extension driving licence
		extensionDl := filepath.Ext(fileDl.Filename)
		// Generate filename driving licence
		var filenameDl string
		filenameDl = lib.RandStringBytesMaskImprSrc(20)
		if fileDl != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			var httpError *echo.HTTPError
			httpError = uploadDrivingLicence(fileDl, filenameDl, extensionDl)
			if httpError.Code != http.StatusOK {
				return httpError
			}
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
		params["driving_licence"] = filenameDl + extensionDl
	}

	var fileAvatar *multipart.FileHeader
	fileAvatar, err = c.FormFile("avatar")
	if fileAvatar != nil {
		// Get file extension driving licence
		extensionAvatar := filepath.Ext(fileAvatar.Filename)
		// Generate filename driving licence
		var filenameAvatar string
		filenameAvatar = lib.RandStringBytesMaskImprSrc(20)
		if fileAvatar != nil {
			if err != nil {
				return lib.CustomError(http.StatusBadRequest)
			}
			var httpError *echo.HTTPError
			httpError = uploadAvatar(fileAvatar, filenameAvatar, extensionAvatar)
			if httpError.Code != http.StatusOK {
				return httpError
			}
		} else {
			return lib.CustomError(http.StatusBadRequest)
		}
		params["avatar"] = filenameAvatar + extensionAvatar
	}
	
	statusUpdate := models.UpdateCustomer(params, idStr)
	if statusUpdate != http.StatusOK {
		return lib.CustomError(http.StatusInternalServerError)
	}
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	
	return c.JSON(http.StatusOK, response)
}