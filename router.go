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
  api := e.Group("/apiv1")
  // Auth
  api.POST("/signup", controllers.SignUp)
  api.POST("/signin", controllers.SignIn)
  api.POST("/logout", controllers.LogOut)
  api.PUT("/resetpassword/:id", controllers.ResetPassword)

  // Customer
  api.GET("/customerdata/:id", controllers.CustomerData)
  api.PUT("/customer/:id", controllers.UpdateCustomerData)

  // Bidding
  api.GET("/biddings", controllers.GetAllBidding)
  api.PUT("/approvebidding/:id", controllers.ApproveBidding)

  // Work Order
  api.POST("/workorder", controllers.CreateWorkder)
  api.GET("/workorders", controllers.GetAllWorkOrder)
  api.GET("/workorderdata/:id", controllers.WorkOrderData)

  // New Order
  api.POST("/neworder", controllers.CreateNewOrder)

  
  

  // Start server
  e.Logger.Fatal(e.Start(":1323"))

  e.Any("/*", func(c echo.Context) error {
	return lib.CustomError(http.StatusMethodNotAllowed)
  })

  return e

}