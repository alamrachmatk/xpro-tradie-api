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
  apps := e.Group("/apiv1")
  // Auth
  e.POST("/signup", controllers.SignUp)
  e.POST("/signin", controllers.SignIn)
  e.POST("/logout", controllers.LogOut)
  e.PUT("/resetpassword/:id", controllers.ResetPassword)

  // Customer
  e.GET("/customerdata/:id", controllers.CustomerData)
  e.PUT("/customer/:id", controllers.UpdateCustomerData)

  // Bidding
  e.GET("/biddings", controllers.GetAllBidding)
  apps.PUT("/approvebidding/:id", controllers.ApproveBidding)

  // Work Order
  apps.POST("/workorder", controllers.CreateWorkder)
  apps.GET("/workorders", controllers.GetAllWorkOrder)
  apps.GET("/workorderdata/:id", controllers.WorkOrderData)

  // New Order
  apps.POST("/neworder", controllers.CreateNewOrder)

  
  

  // Start server
  e.Logger.Fatal(e.Start(":1323"))

  e.Any("/*", func(c echo.Context) error {
	return lib.CustomError(http.StatusMethodNotAllowed)
  })

  return e

}