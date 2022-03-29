package middleware

import (
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func PermInterceptor(perm string) iris.Handler {
	return func(ctx iris.Context) {
		auth, _ := ctx.Values().Get("auth").(*model.AuthInfo)
		role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
		if err := dao.CheckPermission(role, perm); err != nil {
			ctx.JSON(model.ErrorNoPermissions(err))
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
