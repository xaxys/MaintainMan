package middleware

import (
	"fmt"
	"maintainman/config"
	"maintainman/model"

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
	HeaderExtractor = jwt.New(headerJwtConfig).Serve

	TokenValidator = func(ctx iris.Context) {
		jwtToken, ok := ctx.Values().Get("jwt").(*jwt.Token)
		if ok {
			jwtInfo := jwtToken.Claims.(jwt.MapClaims)
			uid := uint(jwtInfo["user_id"].(float64))
			role := jwtInfo["user_role"].(string)
			auth := &model.AuthInfo{
				User: uid,
				Role: role,
				IP:   ctx.Request().RemoteAddr,
			}
			ctx.Values().Set("auth", auth)
		}
		ctx.Next()
	}

	LoginInterceptor = func(ctx iris.Context) {
		if ctx.Values().Get("auth").(*model.AuthInfo) == nil {
			ctx.JSON(model.ErrorUnauthorized(fmt.Errorf("无法获取到登陆用户信息，请重新登陆")))
			ctx.StopExecution()
		}
		ctx.Next()
	}
}

var headerJwtConfig = jwt.Config{
	ErrorHandler: func(ctx iris.Context, err error) {
		if err != nil {
			ctx.JSON(model.ErrorUnauthorized(err))
			ctx.StopExecution()
		}
	},

	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	},

	SigningMethod: jwt.SigningMethodHS256,
	Extractor:     jwt.FromAuthHeader,
}
