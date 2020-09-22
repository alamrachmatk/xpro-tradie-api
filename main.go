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
  // Topics
  e.GET("/topics", controllers.GetAllTopics)
  e.POST("/topics", controllers.CreateTopic)
  e.PUT("/topics/:id", controllers.UpdateTopic)
  e.DELETE("/topics/:id", controllers.DeleteTopic)
  // News
  e.GET("/news", controllers.GetAllNews)
  e.POST("/news", controllers.CreateNews)
  e.PUT("/news/:id", controllers.UpdateNews)
  e.DELETE("/news/:id", controllers.DeleteNews)
  // Customers
  e.POST("/signup", controllers.SignUp)

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}