package main

import (
	"api/controllers"
	"api/db"
	"api/lib"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter() *echo.Echo {
	// Initialize main database
	db.Db = db.MariaDBInit()
	db.RedisPool = db.RedisPoolInit()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:9527"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	api := e.Group("/apiv1")
	api.GET("/dns", controllers.GetAllDns)
	api.GET("/sites", controllers.GetAllSites)
	api.GET("/totaldns", controllers.GetTotalDns)
	api.GET("/totalblock", controllers.GetTotalBlock)
	api.GET("/totaldnsblock", controllers.GetTotalDnsBlock)
	api.GET("/totalipaddress", controllers.GetTotalIpAddress)
	api.GET("/totaltopmostactivelist", controllers.GeTotalTopMostActiveList)
	api.GET("/totaldnsdaylist", controllers.GeTotalDnsDayList)
	api.GET("/totalipaddressdaylist", controllers.GeTotalIpAddressDayList)
	api.GET("/totalrequestlist", controllers.GeTotalRequestList)
	//api.GET("/totaldnsblockcategorydaylist", controllers.GeTotalDnsBlockCategoryDayList)
	//api.GET("/totalipaddressblockcategorydaylist", controllers.GeTotalIpAddressBlockCategoryDayList)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	e.Any("/*", func(c echo.Context) error {
		return lib.CustomError(http.StatusMethodNotAllowed)
	})

	return e

}
