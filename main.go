package main

import (
	controller "example.com/Instahyre/teleapi/controllers"
	ini "example.com/Instahyre/teleapi/init"
	"example.com/Instahyre/teleapi/middleware"
	"github.com/gin-gonic/gin"
)

func init() {

	ini.ConnectToDB()
}

func main() {

	r := gin.Default()
	//r.POST("/create", controller.CreateUser)
	r.POST("/signup", controller.SignUp)
	r.POST("/login", controller.Login)
	r.GET("/search-by-name/:name", middleware.RequireAuth, controller.Validate, controller.SearchByName)
	r.GET("/search-by-number/:number", middleware.RequireAuth, controller.Validate, controller.SearchByNumber)
	r.PATCH("/report-number/:number", middleware.RequireAuth, controller.Validate, controller.Report)
	//user := r.Group("/api/users")

	r.Run("localhost:3000")
}
