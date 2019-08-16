package routes

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"minerva/src/common"
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

	// 登录相关路由
	registerLogin()
	// 注册user相关路由
	registerUser()

}

/**
  注册user 路由
*/
func registerUser() {
	userController := controller.UserController{}

	// Group自动加prefix:
	userRoutes := Route.Group("users")
	// 暂时先做jwt的token校验
	userRoutes.Use(middleware.JWT([]byte(common.JWTKey)))

	// hello
	userRoutes.GET("/hello", userController.SayHello)
	// user列表
	userRoutes.GET("/index", userController.Index)
	// user详情
	userRoutes.GET("/:id", userController.Detail)

}

/**
  注册login 路由
*/
func registerLogin() {
	loginController := controller.LoginController{}
	Route.POST("/login", loginController.Login)
}
