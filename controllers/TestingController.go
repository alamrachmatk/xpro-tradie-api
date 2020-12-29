package controllers

import (
	"api/config"
	"api/lib"
	"api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/jpillora/go-tld"
	"github.com/labstack/echo"
)

func ParsingDomain(c echo.Context) error {
	urls := []string{
		"https://pull-flv-f11-ab.tiktokcdn.com",
		"https://m25.gmdmpdm.com",
		"https://www.medi-cal.ca.gov/",
		"https://ato.gov.au",
		"http://a.very.complex-domain.co.uk:8080/foo/bar",
		"http://a.domain.that.is.unmanaged",
	}
	for _, url := range urls {
		u, _ := tld.Parse(url)
		log.Println(u.Domain)
		// fmt.Printf("%50s = [ %s ] [ %s ] [ %s ] [ %s ] [ %s ] [ %t ]\n",
		// 	u, u.Subdomain, u.Domain, u.TLD, u.Port, u.Path, u.ICANN)
	}
	return nil
}

func ConvertDate(c echo.Context) error {

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
	var responseData []models.DomainExampleList
	var httpError *echo.HTTPError
	total, responseData, httpError = TestingListQuery(limit, offset, pagination)
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

func TestingListQuery(limit uint64, offset uint64, pagination bool) (uint64, []models.DomainExampleList, *echo.HTTPError) {
	var domainexamples []models.DomainExample
	total, err := models.GetAllDomainExample(&domainexamples, limit, offset, pagination, nil)
	if err != nil {
		return 0, nil, lib.CustomError(http.StatusInternalServerError)
	}
	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	var responseData []models.DomainExampleList
	for _, domainexample := range domainexamples {
		log.Println(*domainexample.DateTime)
	}

	if total == 0 {
		return 0, nil, lib.CustomError(http.StatusNotFound)
	}

	return total, responseData, lib.CustomError(http.StatusOK)
}
