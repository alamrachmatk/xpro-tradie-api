package main

import (
  "api/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"api/controllers"
)

func main() {
  // Initialize main database
  db.Db = db.MariaDBInit()
  db.RedisPool = db.RedisPoolInit()
  
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Routes
  // Customers
  e.POST("/signup", controllers.SignUp)
  e.POST("/signin", controllers.SignIn)

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}