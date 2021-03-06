package middleware

import (
	"fmt"

	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/model"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var (
	HeaderExtractor  iris.Handler
	TokenValidator   iris.Handler
	LoginInterceptor iris.Handler
)

var jwtkey = []byte(config.AppConfig.GetString("token.key"))

func init() {
	HeaderExtractor = jwt.New(jwt.Config{
		SigningMethod:       jwt.SigningMethodHS256,
		Extractor:           jwt.FromAuthHeader,
		CredentialsOptional: true,
		ValidationKeyGetter: func(token *jwt.Token) (any, error) {
			return jwtkey, nil
		},
		ErrorHandler: func(ctx iris.Context, err error) {
			response := model.ErrorUnauthorized(err)
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
			ctx.StopExecution()
		},
	}).Serve

	TokenValidator = func(ctx iris.Context) {
		jwtToken, ok := ctx.Values().Get("jwt").(*jwt.Token)
		if ok {
			jwtInfo := jwtToken.Claims.(jwt.MapClaims)
			uid := uint(jwtInfo["user_id"].(float64))
			name := jwtInfo["user_name"].(string)
			role := jwtInfo["user_role"].(string)
			auth := &model.AuthInfo{
				User:  uid,
				Name:  name,
				Role:  role,
				IP:    ctx.Request().RemoteAddr,
				Other: jwtInfo,
			}
			ctx.Values().Set("auth", auth)
		}
		ctx.Next()
	}

	LoginInterceptor = func(ctx iris.Context) {
		if ctx.Values().Get("auth") == nil {
			response := model.ErrorNoPermissions(fmt.Errorf("权限不足：需要登录"))
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
