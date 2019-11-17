package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"minerva/src/common"
	"minerva/src/logic"
	logicCommon "minerva/src/logic/common"
	"minerva/src/model"
	"net/http"
	"strconv"
)

type UserController struct {
}

var UserLogic logic.UserLogic = logic.UserLogic{}

/**
	hello world - 正常 set session
    GET /minerva/users/hello

	@params ctx echo.Context
	@return json
*/
func (u UserController) SayHello(ctx echo.Context) error {

	sessionData, _ := logicCommon.CookieSession.Get(ctx.Request(), "auth")
	var user string = (sessionData.Values["email"]).(string)
	var result string = "hello " + user
	return ctx.String(http.StatusOK, result)
}

/**
	hello world - 读取token
    GET /minerva/users/helloByToken

	@params ctx echo.Context
	@return json
*/
func (u UserController) SayHelloByToken(ctx echo.Context) error {
	// authorization := ctx.Request().Header.Get("Authorization")
	// authorization:    Bearer ***

	users := ctx.Get("jwt_auth").(*jwt.Token)
	claims := users.Claims.(*JwtCustomClaims)

	common.Logger.Println(claims)
	return ctx.String(http.StatusOK, "Welcome "+claims.Email+"!")
}

/**
	获取user详情
    GET /minerva/users/:id
	@params ctx echo.Context
	@return json
*/
func (u UserController) Detail(ctx echo.Context) error {
	queryId := ctx.Param("id")
	id, error := strconv.Atoi(queryId)
	// 处理 error
	if error != nil {
		common.Logger.Fatal("user #Detail : error: ", error)
	}

	var user *model.User = UserLogic.Detail(id)

	data := make(map[int]model.User)
	data[id] = *user
	return ctx.JSON(http.StatusOK, data)
}

/**
	获取user列表
    GET /minerva/users/index
	@params ctx echo.Context
	@return json
*/
func (u UserController) Index(ctx echo.Context) error {
	// 验证。。
	params := make(map[string]interface{})
	data := UserLogic.Index(params, ctx)

	return ctx.JSON(http.StatusOK, data)
}
