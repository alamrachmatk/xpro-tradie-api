package controllers

import (
	"api/config"
	"api/db"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/garyburd/redigo/redis"
)

func GetAllBidding(c echo.Context) error {
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
	var responseData []models.BiddingList
	var httpError *echo.HTTPError
	total, responseData, httpError = BiddingListQuery(limit, offset, pagination)
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

func BiddingListQuery(limit uint64, offset uint64, pagination bool) (uint64, []models.BiddingList, *echo.HTTPError){
	var biddings []models.Bidding
	total, err := models.GetAllBidding(&biddings, limit, offset, pagination, nil)
	if err != nil {
		return 0, nil, lib.CustomError(http.StatusInternalServerError)
	}
	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	var companyIDs []string
	for _, bidding := range biddings {
		if bidding.CompanyID > 0 {
			companyIDs = append(companyIDs, strconv.FormatUint(bidding.CompanyID, 10))
		}
	}

	var companies []models.Company
	models.GetCompanyIn(&companies, companyIDs...)

	companiesMap := make(map[uint64]models.Company)

	for _, company := range companies {
		companiesMap[company.CompanyID] = company
	}

	var responseData []models.BiddingList
	for _, bidding := range biddings {
		var data models.BiddingList
		data.BiddingID = bidding.BiddingID
		data.OrderID = bidding.OrderID
 
		if bidding.CompanyID > 0 {
			company := companiesMap[bidding.CompanyID]
			data.Company.CompanyID = company.CompanyID
			data.Company.Name = company.Name
		}

		data.LaborTime = bidding.LaborTime
		data.Price = bidding.Price
		data.Description = bidding.Description
		data.Status = bidding.Status

		responseData = append(responseData, data)
	}

	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	return total, responseData, lib.CustomError(http.StatusOK)
}

func ApproveBidding(c echo.Context) error {
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

	idStr := c.Param ("id")
	if idStr != "" {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if id > 0 {
			var bidding models.Bidding
			status := models.GetBidding(&bidding, idStr)
			if status != http.StatusOK {
				return lib.CustomError(http.StatusNotFound)
			}
		}
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	params := make(map[string]string)

	// Get parameter status
	status := c.FormValue("status")
	if status != "" {
		params["status"] = status
	} else {
		return lib.CustomError(http.StatusBadRequest)
	}

	statusUpdate := models.UpdateBidding(params, idStr)
	if statusUpdate != http.StatusOK {
		return lib.CustomError(http.StatusInternalServerError)
	}

	var response lib.Response
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"

	return c.JSON(http.StatusOK, response)


}