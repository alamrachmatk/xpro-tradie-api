package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"
	"github.com/labstack/echo"
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