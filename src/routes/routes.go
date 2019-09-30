package routes

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"minerva/src/common"
	"minerva/src/http/controller"
)

// 路由管理
var Route *echo.Echo

// base 路由组
var BaseGroup *echo.Group

/**
  注册路由
  对外暴露的方法
*/
func RegisterRoutes() {
	// echo的路由
	Route = echo.New()

	// 服务名作 基础路由的group
	var serviceName string = viper.GetString("common.serviceName")
	BaseGroup = Route.Group(serviceName)

	// recover异常
	BaseGroup.Use(middleware.Recover())

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
	userRoutes := BaseGroup.Group("/users")

	// jwt config
	config := middleware.JWTConfig{
		Claims:     &controller.JwtCustomClaims{},
		SigningKey: []byte(common.JWTKey),
		// 获取token的标志
		ContextKey: "jwt_auth",
	}
	// jwt 使用
	userRoutes.Use(middleware.JWTWithConfig(config))

	// 暂时先做jwt的token校验
	// userRoutes.Use(middleware.JWT([]byte(common.JWTKey)))

	// hello world
	userRoutes.GET("/hello", userController.SayHello)
	// hello2
	userRoutes.GET("/helloByToken", userController.SayHelloByToken)

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

	BaseGroup.POST("/login", loginController.Login)
}
