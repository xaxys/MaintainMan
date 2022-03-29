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
			response := model.ErrorNoPermissions(err)
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
