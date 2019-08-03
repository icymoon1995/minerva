package controller

import (
	"github.com/labstack/echo"
	"log"
	"minerva/src/logic"
	"minerva/src/model"
	"net/http"
	"strconv"
)

type UserController struct {
}

var UserLogic logic.UserLogic = logic.UserLogic{}

// hello world
func (u UserController) SayHello(ctx echo.Context) error {
	var result string = "hello world"
	return ctx.String(http.StatusOK, result)
}

/**
	获取user详情
    GET /users/:id
	@params ctx echo.Context
	@return json
*/
func (u UserController) Detail(ctx echo.Context) error {
	// formId := ctx.FormValue("id")
	queryId := ctx.Param("id")
	id, error := strconv.Atoi(queryId)
	// 处理 error
	if error != nil {
		log.Fatal(error)
	}

	var user *model.User = UserLogic.Detail(id)

	data := make(map[int]model.User)
	data[id] = *user
	return ctx.JSON(http.StatusOK, data)
}

/**
	获取user列表
    GET /users/:id
	@params ctx echo.Context
	@return json
*/
func (u UserController) Index(ctx echo.Context) error {
	// 验证。。
	params := make(map[string]interface{})
	data := UserLogic.Index(params)

	return ctx.JSON(http.StatusOK, data)
}
