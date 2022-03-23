package controller

import (
	"maintainman/model"
	"maintainman/service"
	"maintainman/util"

	"github.com/kataras/iris/v12"
)

func GetUserOrders(ctx iris.Context) {
	id, _ := ctx.Values().GetUint("user_id")
	status, _ := ctx.URLParamInt("status")
	offset, _ := ctx.URLParamInt("offset")
	response := service.GetOrderByUser(id, util.ToUint(status), util.ToUint(offset))
	ctx.Values().Set("response", response)
}

func GetRepairerOrders(ctx iris.Context) {
	id, _ := ctx.Values().GetUint("user_id")
	current, _ := ctx.URLParamBool("current")
	offset, _ := ctx.URLParamInt("offset")
	response := service.GetOrderByRepairer(id, current, util.ToUint(offset))
	ctx.Values().Set("response", response)
}

func GetAllOrders(ctx iris.Context) {
	aul := &model.AllOrderJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.GetAllOrders(aul)
	ctx.Values().Set("response", response)
}

func GetOrderByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	response := service.GetOrderByID(id)
	ctx.Values().Set("response", response)
}

func CreateOrder(ctx iris.Context) {
	aul := &model.ModifyOrderJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	response := service.CreateOrder(aul)
	ctx.Values().Set("response", response)
}

func UpdateOrder(ctx iris.Context) {
	aul := &model.ModifyOrderJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	id, _ := ctx.Params().GetUint("id")
	response := service.UpdateOrder(id, aul)
	ctx.Values().Set("response", response)
}

func UpdateOrderByID(ctx iris.Context) {
	aul := &model.ModifyOrderJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	id, _ := ctx.Params().GetUint("id")
	response := service.UpdateOrder(id, aul)
	ctx.Values().Set("response", response)
}

// change order status

func ReleaseOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.ReleaseOrder(id, uid)
	ctx.Values().Set("response", response)
}

func AssignOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	repairer, _ := ctx.URLParamInt("repairer")
	response := service.AssignOrder(id, uid, util.ToUint(repairer))
	ctx.Values().Set("response", response)
}

func SelfAssignOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.SelfAssignOrder(id, uid)
	ctx.Values().Set("response", response)
}

func CompleteOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.CompleteOrder(id, uid)
	ctx.Values().Set("response", response)
}

func CancelOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.CancelOrder(id, uid)
	ctx.Values().Set("response", response)
}

func RejectOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.RejectOrder(id, uid)
	ctx.Values().Set("response", response)
}

func ReportOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.ReportOrder(id, uid)
	ctx.Values().Set("response", response)
}

func HoldOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	response := service.HoldOrder(id, uid)
	ctx.Values().Set("response", response)
}

func AppraiseOrder(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	uid, _ := ctx.Values().GetUint("user_id")
	appraisal, _ := ctx.URLParamInt("appraisal")
	response := service.AppraiseOrder(id, util.ToUint(appraisal), uid)
	ctx.Values().Set("response", response)
}
