package middleware

import (
	"maintainman/dao"
	"maintainman/model"

	"github.com/kataras/iris/v12"
)

func PermInterceptor(perm string) iris.Handler {
	return func(ctx iris.Context) {
		role := ""
		if auth, ok := ctx.Values().Get("auth").(*model.AuthInfo); ok {
			role = auth.Role
		}
		if err := dao.CheckPermission(role, perm); err != nil {
			ctx.JSON(model.ErrorNoPermissions(err))
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
