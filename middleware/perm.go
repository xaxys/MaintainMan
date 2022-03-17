package middleware

import (
	"maintainman/dao"
	"maintainman/model"

	"github.com/kataras/iris/v12"
)

func PermInterceptor(perm string) iris.Handler {
	return func(ctx iris.Context) {
		urole := ctx.Values().GetString("user_role")
		if err := dao.CheckPermission(urole, perm); err != nil {
			ctx.JSON(model.ErrorNoPermissions(err))
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
