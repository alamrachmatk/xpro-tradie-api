package controllers

import (
	"api/lib"
	"log"
	"api/config"
	"api/models"
	"net/http"
	"strconv"
	"github.com/labstack/echo"
)

func GetAllNews(c echo.Context) error {
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

	statusStr := c.QueryParam("status")

	var total uint64
	var response lib.Response
	var news []models.News
	var topicIDs []string
	var topic []models.Topic
	var data models.NewsList
	var responseData []models.NewsList
	var i uint64
	totalTopic, err := models.GetAllTopics(&topic, limit, offset, pagination, nil)
	log.Println(totalTopic)
	for _, topics := range topic {
		topicIDs = append(topicIDs, strconv.FormatUint(topics.TopicID, 10))
	} 
	log.Println(statusStr)
	models.APIGetNewsIn(&news, topicIDs...)
	for _, dnews := range news {
		for _, topics := range topic {
			if topics.TopicID == dnews.TopicID {
				if statusStr == dnews.Status {
					i++
					data.NewsID = dnews.NewsID
					data.Topic 	= topics.Name
					data.News 	= dnews.News
					data.Description = dnews.Description
					data.Status	= dnews.Status
					data.DateCreated = dnews.DateCreated
					responseData = append(responseData, data)
					total = i
				} else if statusStr == "" {
					if dnews.Status == "publish" {
						i++
						data.NewsID = dnews.NewsID
						data.Topic 	= topics.Name
						data.News 	= dnews.News
						data.Description = dnews.Description
						data.Status	= dnews.Status
						data.DateCreated = dnews.DateCreated
						responseData = append(responseData, data)
						total = i
					}
				}
			}
		}
	}
	  
	if total == 0 {
		return lib.CustomError(http.StatusNotFound, "Not found")
	}
 
	meta := lib.GenerateMeta(c, total, limit, page, offset, pagination)
	response.Status.Code = http.StatusOK
	response.Status.MessageServer = "OK"
	response.Status.MessageClient = "OK"
	response.Meta = &meta
	response.Data = responseData
	return c.JSON(http.StatusOK, response)

}

func CreateNews(c echo.Context) error {
	var news map[string]string
	news = make(map[string]string)
	// Get parameter name
	topicID := c.FormValue("topic_id")
	if topicID != "" {
		news["topic_id"] = topicID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Get parameter name
	name := c.FormValue("name")
	if name != "" {
		news["name"] = name
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Get parameter description
	description := c.FormValue("description")
	if description != "" {
		news["description"] = description
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Set parameter status
	news["status"] = "draft"
	status := models.CreateNews(news)
	return echo.NewHTTPError(status)
}

func UpdateNews(c echo.Context) error {
	params := make(map[string]string)
	// Get parameter id
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if id > 0 {
		params["id"] = strconv.FormatUint(id, 10)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Get parameter id
	topicIDStr := c.FormValue("topic_id")
	topicID, _ := strconv.ParseUint(topicIDStr, 10, 64)
	log.Println("topicID", topicID)
	if topicID > 0 {
		params["topic_id"] = strconv.FormatUint(topicID, 10)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// Get parameter name
	name := c.FormValue("name")
	if name != "" {
		params["name"] = name
	}
	// Get parameter name
	description := c.FormValue("description")
	if description != "" {
		params["description"] = description
	}
	// Get parameter status
	statusp := c.FormValue("status")
	if statusp != "" {
		params["status"] = statusp
	}

	status := models.UpdateNews(params)
	return echo.NewHTTPError(status)
}

func DeleteNews(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	status := models.DeleteNews(strconv.FormatUint(id, 10))
	return echo.NewHTTPError(status)
}