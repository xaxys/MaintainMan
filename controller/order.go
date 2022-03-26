package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetUserOrders(ctx iris.Context) {
	req := &model.UserOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByUser(req, auth)
	ctx.Values().Set("response", response)
}

func GetRepairerOrders(ctx iris.Context) {
	req := &model.RepairerOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByRepairer(auth.User, req, auth)
	ctx.Values().Set("response", response)
}

func ForceGetRepairerOrders(ctx iris.Context) {
	req := &model.RepairerOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByRepairer(id, req, auth)
	ctx.Values().Set("response", response)
}

func GetAllOrders(ctx iris.Context) {
	req := &model.AllOrderRequest{}
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetAllOrders(req, auth)
	ctx.Values().Set("response", response)
}

func GetOrderByID(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.GetOrderByID(id, auth)
	ctx.Values().Set("response", response)
}

func CreateOrder(ctx iris.Context) {
	aul := &model.CreateOrderRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CreateOrder(aul, auth)
	ctx.Values().Set("response", response)
}

func UpdateOrder(ctx iris.Context) {
	aul := &model.UpdateOrderRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.UpdateOrder(id, aul, auth)
	ctx.Values().Set("response", response)
}

func ForceUpdateOrder(ctx iris.Context) {
	aul := &model.UpdateOrderRequest{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ForceUpdateOrder(id, aul, auth)
	ctx.Values().Set("response", response)
}

// change order status

func ReleaseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ReleaseOrder(id, auth)
	ctx.Values().Set("response", response)
}

func AssignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	repairer := util.ToUint(ctx.URLParamIntDefault("repairer", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AssignOrder(id, repairer, auth)
	ctx.Values().Set("response", response)
}

func SelfAssignOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AssignOrder(id, auth.User, auth)
	ctx.Values().Set("response", response)
}

func CompleteOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CompleteOrder(id, auth)
	ctx.Values().Set("response", response)
}

func CancelOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.CancelOrder(id, auth)
	ctx.Values().Set("response", response)
}

func RejectOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.RejectOrder(id, auth)
	ctx.Values().Set("response", response)
}

func ReportOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.ReportOrder(id, auth)
	ctx.Values().Set("response", response)
}

func HoldOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.HoldOrder(id, auth)
	ctx.Values().Set("response", response)
}

func AppraiseOrder(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	appraisal := util.ToUint(ctx.URLParamIntDefault("appraisal", 0))
	auth := util.NilOrPtrCast[model.AuthInfo](ctx.Values().Get("auth"))
	response := service.AppraiseOrder(id, appraisal, auth)
	ctx.Values().Set("response", response)
}
