package routes

import (
	"github.com/labstack/echo"
	"minerva/src/http/controller"
)

var Route *echo.Echo

/**
  注册路由
  对外暴露的方法
 */
func RegisterRoutes() {
	// echo的路由
	Route = echo.New()

	// 注册user相关路由
	registerUser()

}

/**
  注册user 路由
 */
func registerUser() {
	userController := controller.UserController{}

	// Group自动加prefix:
	userRoutes := Route.Group("user")

	userRoutes.GET("/hello", userController.SayHello)
	userRoutes.GET("/user", userController.Get2)
	userRoutes.GET("/index", userController.GetUser)
}