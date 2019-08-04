package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"minerva/src/common"
	"minerva/src/logic"
	"net/http"
	"strconv"
	"time"
)

type LoginController struct {
}

var authLogic logic.AuthLogic = logic.AuthLogic{}

// jwt加密的数据结构体
type JwtCustomClaims struct {
	Email string `json:"email"`
	IsGod bool   `json:"admin"`
	jwt.StandardClaims
}

func (login LoginController) Login(ctx echo.Context) error {

	// 用户名
	email := ctx.FormValue("email")
	// 密码
	password := ctx.FormValue("password")

	// 验证用户名密码
	error := verify(email, password)
	if error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, error.Error())
	}

	// 生成token的必要数据
	data := map[string]interface{}{
		"email": email,
		"isGod": true,
	}

	// 生成token
	token, error := generateToken(data)

	if error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, error.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token,
	})

	return echo.ErrUnauthorized
}

/**
校验用户名和密码
*/
func verify(email string, password string) error {
	// 后续考虑第三方认证。。。。todo
	result, error := authLogic.Verify(email, password)
	if !result {
		return error
	}

	return nil
}

/**
利用JWT生成token
*/
func generateToken(data map[string]interface{}) (string, error) {
	var env string = common.Enviorment + "."
	// jwt private_key
	var jwtKey string = common.JWTKey
	// 过期时间 -- 必须转化成time.Duration格式 不然会抛异常
	var jwtExpire string = viper.GetString(env + "jwt.expire")
	jwtExpireInt, _ := strconv.ParseInt(jwtExpire, 10, 64)
	var d time.Duration = time.Duration(jwtExpireInt) * time.Hour

	// 填充必要数据 目前只用了email 和 isGod
	claims := &JwtCustomClaims{
		data["email"].(string),
		data["isGod"].(bool),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(d).Unix(),
		},
	}

	//	通过claims生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// jwt签名
	reallyToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return reallyToken, nil
}
