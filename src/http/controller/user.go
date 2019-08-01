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

// controller -> logic -> model 连贯
// 返回一个传统的json
func (u UserController) GetUser(ctx echo.Context) error {
	formId := ctx.FormValue("id")
	id, error := strconv.Atoi(formId)
	// 处理 error
	if error != nil {
		log.Fatal(error)
	}

	var user *model.User = UserLogic.Detail(id)

	data := make(map[int]interface{})
	data[id] = user
	// return ctx.String(http.StatusOK, user.Name)
	return ctx.JSON(http.StatusOK, data)
}