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
	api.GET("/totaldns", controllers.GetTotalDns)
	api.GET("/totalblok", controllers.GetTotalBlok)
	api.GET("/totaldnsblok", controllers.GetTotalDnsBlok)
	api.GET("/totaltopmostactivelist", controllers.GeTotalTopMostActiveList)
	api.GET("/dns", controllers.GetAllDns)
	api.GET("/sites", controllers.GetAllSites)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	e.Any("/*", func(c echo.Context) error {
		return lib.CustomError(http.StatusMethodNotAllowed)
	})

	return e

}
