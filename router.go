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
	api.GET("/totalsite", controllers.GetTotalSite)
	api.GET("/totaltopmostactivelist", controllers.GeTotalTopMostActiveList)
	api.GET("/sites", controllers.GetAllSites)
	api.GET("/parsingdomain", controllers.ParsingDomain)
	api.GET("/convertdate", controllers.ConvertDate)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	e.Any("/*", func(c echo.Context) error {
		return lib.CustomError(http.StatusMethodNotAllowed)
	})

	return e

}
