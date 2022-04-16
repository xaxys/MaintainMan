package rbac

import (
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"github.com/kataras/iris/v12"
)

func PermInterceptor(perm string) iris.Handler {
	if PermPO.perm[perm] == "" {
		logger.Logger.Errorf("Permission not declared: %s", perm)
	}
	logger.Logger.Debugf("Permission Registered: %s", perm)
	return func(ctx iris.Context) {
		auth, _ := ctx.Values().Get("auth").(*model.AuthInfo)
		role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
		if err := CheckPermission(role, perm); err != nil {
			response := model.ErrorNoPermissions(err)
			ctx.StatusCode(response.Code)
			ctx.JSON(response)
			ctx.StopExecution()
		}
		ctx.Next()
	}
}
