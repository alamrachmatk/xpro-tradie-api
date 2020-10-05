package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"encoding/base64"
	"log"	
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
	"time"


	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/disintegration/imaging"
	jwt "github.com/dgrijalva/jwt-go"
)

func SignUp(c echo.Context) error  {
	var err error
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

	// Get parameter password
	password := c.FormValue("password")
	if password != "" {
		passwordEncode := base64.StdEncoding.EncodeToString([]byte(password))
		params["password"] = passwordEncode
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

	// Get parameter photo id
	photoId := c.FormValue("photo_id")
	if photoId != "" {
		params["photo_id"] = photoId
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	params["status"] = "1"

	var fileDl *multipart.FileHeader
	fileDl, err = c.FormFile("driving_licence")
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

	var fileAvatar *multipart.FileHeader
	fileAvatar, err = c.FormFile("avatar")
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

	var customer models.Customer
	status := models.GetCustomerEmail(&customer, email)
	if status == http.StatusOK {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Forbidden, email address is already registered")
	}

	statusResponse, lastID := models.CreateCustomer(params)
	params["customer_id"] = strconv.Itoa(int(lastID))
	if statusResponse != 200 {
		return lib.CustomError(http.StatusInternalServerError)
	}
	
	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)
}

func uploadDrivingLicence(file *multipart.FileHeader, filenameDl string, extension string)  (*echo.HTTPError) {
		var err error
		// Upload image and move to proper directory
		err = lib.UploadImage(file, config.MediaServerPathDl+"/"+filenameDl+extension)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		// Open a test image.
		src, err := imaging.Open(config.MediaServerPathDl+"/"+filenameDl+extension)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		src = imaging.Resize(src, 500, 0, imaging.Lanczos)
		err = imaging.Save(src, config.MediaServerPathDl+"/"+filenameDl+extension)
		if err != nil {
			log.Println(err)
			return lib.CustomError(http.StatusInternalServerError)
		}
		return echo.NewHTTPError(http.StatusOK)
}

func uploadAvatar(file *multipart.FileHeader, filenameAvatar string, extension string)  (*echo.HTTPError) {
	var err error
	// Upload image and move to proper directory
	err = lib.UploadImage(file, config.MediaServerPathAvatar+"/"+filenameAvatar+extension)
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusInternalServerError)
	}
	// Open a test image.
	src, err := imaging.Open(config.MediaServerPathAvatar+"/"+filenameAvatar+extension)
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusInternalServerError)
	}
	src = imaging.Resize(src, 500, 0, imaging.Lanczos)
	err = imaging.Save(src, config.MediaServerPathAvatar+"/"+filenameAvatar+extension)
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusInternalServerError)
	}
	return echo.NewHTTPError(http.StatusOK)
}

func SignIn(c echo.Context) error {
	var err error
	params:= make(map[string]string)
	request := c.Request()
	useragent := request.Header.Get("User-Agent")

	// Get parameter email
	email := c.FormValue("email")
	if email != "" {
		params["email"] = email
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	// Get parameter password
	password := c.FormValue("password")
	if email != "" {
		params["password"] = password
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	var customer models.Customer
	checkEmail := models.GetCustomerEmail(&customer, email)
	if checkEmail != http.StatusOK {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Forbidden, email is not registered")
	}
	
	passwordDecode, _ := base64.StdEncoding.DecodeString(customer.Password)
	if password != string(passwordDecode) {
		return lib.CustomError(http.StatusUnauthorized, "Unauthorized", "password invalid")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customer.CustomerID,
		"email": customer.Email,
		"password":  customer.Password,
		"user-agent": useragent,
		"ts":    time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		log.Println(err)
		return lib.CustomError(http.StatusBadGateway, "error SignedString")
	}

	redisPoolToken := db.RedisPool.Get()
	defer redisPoolToken.Close()
	redisPoolToken.Do("SELECT", config.RedisDBCacheToken)
	redisPoolToken.Do("SET", tokenString, "true")
	redisPoolToken.Do("EXPIRE", tokenString, config.ExpireToken)

	var response lib.ResponseToken
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Email = customer.Email
	response.Token = tokenString

	return c.JSON(http.StatusOK, response)
}

func LogOut(c echo.Context) error {
	// Get authorization token
	var token string
	request := c.Request()
	authorization := request.Header["Authorization"]
	if authorization != nil {
		if strings.HasPrefix(authorization[0], "Bearer ") == true {
			token = authorization[0][7:]
			log.Println(token)
		}
	}

	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	redisPool.Do("SELECT", config.RedisDBCacheToken)
	_, err := redis.Bytes(redisPool.Do("GET", token))
	if err != nil {
		log.Println("No user login with browser:", token)
		return lib.CustomError(http.StatusUnauthorized, "Unauthorized")
	}

	// Remove token from Redis
	redisPoolRt := db.RedisPool.Get()
	defer redisPoolRt.Close()
	redisPool.Do("SELECT", config.RedisDBCacheToken)
	redisPool.Do("DEL", token)

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)
}

func ResetPassword(c echo.Context) error {
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
	idStr := c.Param("id")

	var customer models.Customer
	if idStr != "" {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if id > 0 {
			params["id"] = strconv.FormatUint(id,10)
			status := models.GetCustomer(&customer, idStr)
			if status != http.StatusOK {
				return lib.CustomError(http.StatusNotFound)
			}
		}
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}
	 
	oldPassword := c.FormValue("old_password")
	if oldPassword != "" {
		oldPassword = base64.StdEncoding.EncodeToString([]byte(oldPassword))
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	password := c.FormValue("password")
	if password != "" {
		password = base64.StdEncoding.EncodeToString([]byte(password))
		params["password"] = password
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	confirmPassword := c.FormValue("confirm_password")
	if confirmPassword != "" {
		confirmPassword = base64.StdEncoding.EncodeToString([]byte(confirmPassword))
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	if oldPassword != customer.Password {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Forbidden, old password not valid")
	}

	if password == oldPassword {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Forbidden, password already in use")
	}

	if password != confirmPassword {
		return lib.CustomError(http.StatusForbidden, "Forbidden", "Forbidden, password confirm not valid")
	}

	status := models.UpdateCustomer(params, idStr)
	if status != http.StatusOK {
		return lib.CustomError(http.StatusInternalServerError)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	
	return c.JSON(http.StatusOK, response)
}