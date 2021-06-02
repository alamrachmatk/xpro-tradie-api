package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetAllDns(c echo.Context) error {
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

	params := make(map[string]string)

	title := c.QueryParam("title")
	if title != "" {
		params["title"] = title
	}

	var total uint64
	var responseData []models.DnsList
	var httpError *echo.HTTPError
	total, responseData, httpError = DnsListQuery(limit, offset, pagination, params)
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

func DnsListQuery(limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, []models.DnsList, *echo.HTTPError) {
	var dns []models.Dns
	total, err := models.GetAllDns(&dns, limit, offset, pagination, params)
	if err != nil {
		return 0, nil, lib.CustomError(http.StatusInternalServerError)
	}
	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	var responseData []models.DnsList
	for _, d := range dns {
		var data models.DnsList
		data.DnsID = d.DnsID
		data.Domain = d.Domain
		data.BaseDomain = d.BaseDomain
		data.IpAddress = d.IpAddress
		if d.HasSubdomain == 0 {
			data.HasSubdomain = "Not Active"
		} else if d.HasSubdomain == 1 {
			data.HasSubdomain = "Active"
		}
		data.LogDatetime = d.LogDatetime
		data.CreatedAt = d.CreatedAt

		responseData = append(responseData, data)
	}

	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	return total, responseData, lib.CustomError(http.StatusOK)
}
