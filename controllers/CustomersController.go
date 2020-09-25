package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"encoding/json"
	"encoding/base64"
	"log"	
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"


	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/disintegration/imaging"
	jwt "github.com/dgrijalva/jwt-go"
)

func SignUp(c echo.Context) error  {
	params := make(map[string]string)
	var err error
	redisPool := db.RedisPool.Get()
	defer redisPool.Close()

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
		// check duplicate email/telp
		redisPool.Do("SELECT", config.RedisDBCacheCustomerByEmail)
		_, err = redis.Bytes(redisPool.Do("GET", email))
		log.Println("err email" + email)
		if err == nil {
			return lib.CustomError(http.StatusForbidden)
		}
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

	// Get parameter avatar
	status := c.FormValue("status")
	if status != "" {
		params["status"] = status
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

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

	redisPool.Do("SELECT", config.RedisDBCacheCustomerByEmail)
	_, err = redis.Bytes(redisPool.Do("GET", email))
	if err != nil {
		statusResponse, lastID := models.CreateCustomer(params)
		params["id"] = strconv.Itoa(int(lastID))
		if statusResponse != 200 {
			return lib.CustomError(http.StatusInternalServerError)
		}
		dataJSON, _ := json.Marshal(params)
		redisPool.Do("SET", email, dataJSON)
	} else {
		return lib.CustomError(http.StatusForbidden)
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

	params:= make(map[string]string)
	var err error

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

	redisPool := db.RedisPool.Get()
	defer redisPool.Close()
	redisPool.Do("SELECT", config.RedisDBCacheCustomerByEmail)
	dataStr, err := redis.Bytes(redisPool.Do("GET", email))
	if err != nil {
		log.Println("No user with email:", email)
		return lib.CustomError(http.StatusUnauthorized, "Unauthorized", "Please check your email and password")
	}

	var customer models.CustomerData
	err = json.Unmarshal([]byte(dataStr), &customer)
	if err != nil {
		log.Println("error unmarshalProperties:", err)
		return lib.CustomError(http.StatusInternalServerError,
			"Sorry, we've experience a problem. Please try again.",
			"Internal server error")
	}
	
	passwordDecode, _ := base64.StdEncoding.DecodeString(customer.Password)
	if password != string(passwordDecode) {
		log.Println("Invalid password")
		return lib.CustomError(http.StatusUnauthorized, "Unauthorized", "Please check your email and password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
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
	redisPoolToken.Do("SET", tokenString+".exist", "true")
	redisPoolToken.Do("EXPIRE", tokenString+".exist", config.ExpireToken)

	var response lib.ResponseToken
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Token = tokenString
	response.Email = customer.Email
	response.Expire = 232424

	return c.JSON(http.StatusOK, response)
}