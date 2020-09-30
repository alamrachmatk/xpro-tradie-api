package main

import (
  	"api/db"
	"api/controllers"
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

  // Routes
  // Auth
  e.POST("/signup", controllers.SignUp)
  e.POST("/signin", controllers.SignIn)
  e.POST("/logout", controllers.LogOut)

  // Customer
  e.GET("/customerdata", controllers.CustomerData)
  e.PUT("/customer/:id", controllers.UpdateCustomerData)
  // Start server
  e.Logger.Fatal(e.Start(":1323"))

  e.Any("/*", func(c echo.Context) error {
	return lib.CustomError(http.StatusMethodNotAllowed)
  })

  return e

}