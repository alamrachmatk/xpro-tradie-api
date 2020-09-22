package controllers

import (
	"api/lib"
	"api/config"
	"api/models"
	"net/http"
	"strconv"
	"github.com/labstack/echo"
)

func GetAllTopics(c echo.Context) error {
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
			return echo.NewHTTPError(http.StatusBadRequest)
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
			return echo.NewHTTPError(http.StatusBadRequest)
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
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}

	var topics []models.Topic
	total, err := models.GetAllTopics(&topics, limit, offset, pagination, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if total == 0 {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	result := lib.Paginate(c, topics, total, limit, page, offset, pagination)

	return c.JSON(http.StatusOK, result)
}

func CreateTopic(c echo.Context) error {
	var topic map[string]string
	topic = make(map[string]string)
	// Get parameter name
	name := c.FormValue("name")
	if name != "" {
		topic["name"] = name
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Set parameter date_created
	topic["status"] = "publish"
	status := models.CreateTopic(topic)
	return echo.NewHTTPError(status)
}

func UpdateTopic(c echo.Context) error {
	params := make(map[string]string)
	// Get parameter id
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if id > 0 {
		params["id"] = strconv.FormatUint(id, 10)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Get parameter name
	name := c.FormValue("name")
	if name != "" {
		params["name"] = name
	}
	// Get parameter status
	statusp := c.FormValue("status")
	if statusp != "" {
		params["status"] = statusp
	}

	status := models.UpdateTopic(params)
	return echo.NewHTTPError(status)
}

func DeleteTopic(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	status := models.DeleteTopic(strconv.FormatUint(id, 10))
	return echo.NewHTTPError(status)
}
