package order

import (
	"github.com/xaxys/maintainman/core/middleware"
	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/rbac"

	"github.com/kataras/iris/v12"
)

var Module module.Module

func init() {
	Module = module.Module{
		ModuleName:    "order",
		ModuleVersion: "1.1.0",
		ModuleConfig:  orderConfig,
		ModuleEnv: map[string]any{
			"orm.model": []any{
				&Tag{},
				&Order{},
				&Status{},
				&Comment{},
				&Item{},
				&ItemLog{},
			},
		},
		ModuleExport: map[string]any{
			"wechat.status.tmpl":     "",
			"wechat.status.order":    "",
			"wechat.status.title":    "",
			"wechat.status.status":   "",
			"wechat.status.time":     "",
			"wechat.status.other":    "",
			"wechat.comment.tmpl":    "",
			"wechat.comment.title":   "",
			"wechat.comment.name":    "",
			"wechat.comment.messgae": "",
			"wechat.comment.time":    "",
		},
		ModulePerm: map[string]string{
			"order.view":        "查看我的订单",
			"order.viewfix":     "查看我维修的订单",
			"order.create":      "创建订单",
			"order.cancel":      "取消订单",
			"order.update":      "更新订单",
			"order.updateall":   "更新所有订单",
			"order.assign":      "分配订单",
			"order.selfassign":  "给自己分配订单",
			"order.release":     "释放订单",
			"order.reject":      "拒绝订单",
			"order.report":      "上报订单",
			"order.hold":        "挂起订单",
			"order.complete":    "完成订单",
			"order.appraise":    "评价订单",
			"order.viewall":     "查看所有订单",
			"comment.view":      "查看我的评论",
			"comment.create":    "创建评论",
			"comment.delete":    "删除评论",
			"comment.viewall":   "查看所有评论",
			"comment.createall": "创建所有评论",
			"comment.deleteall": "删除所有评论",
			"tag.create":        "创建标签",
			"tag.delete":        "删除标签",
			"tag.view":          "查看标签",
			"tag.add":           "添加标签",
			"item.create":       "创建零件",
			"item.delete":       "删除零件",
			"item.viewall":      "查看所有零件",
			"item.update":       "更新零件",
			"item.consume":      "消耗零件",
		},
		EntryPoint: entry,
	}
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx

	Module.ModuleExport["wechat.status.tmpl"] = orderConfig.GetString("notify.wechat.status.tmpl")
	Module.ModuleExport["wechat.status.order"] = orderConfig.GetString("notify.wechat.status.order")
	Module.ModuleExport["wechat.status.title"] = orderConfig.GetString("notify.wechat.status.title")
	Module.ModuleExport["wechat.status.status"] = orderConfig.GetString("notify.wechat.status.status")
	Module.ModuleExport["wechat.status.time"] = orderConfig.GetString("notify.wechat.status.time")
	Module.ModuleExport["wechat.status.other"] = orderConfig.GetString("notify.wechat.status.other")

	Module.ModuleExport["wechat.comment.tmpl"] = orderConfig.GetString("notify.wechat.comment.tmpl")
	Module.ModuleExport["wechat.comment.title"] = orderConfig.GetString("notify.wechat.comment.title")
	Module.ModuleExport["wechat.comment.name"] = orderConfig.GetString("notify.wechat.comment.name")
	Module.ModuleExport["wechat.comment.message"] = orderConfig.GetString("notify.wechat.comment.message")
	Module.ModuleExport["wechat.comment.time"] = orderConfig.GetString("notify.wechat.comment.time")

	mctx.Scheduler.Every(orderConfig.GetString("appraise.purge")).SingletonMode().Do(autoAppraiseOrderService)

	mctx.Route.Get("/wxtmpl/status", getWxStatusTemplateID)
	mctx.Route.Get("/wxtmpl/comment", getWxCommentTemplateID)

	mctx.Route.PartyFunc("/order", func(order iris.Party) {
		order.Get("/user", rbac.PermInterceptor("order.view"), getUserOrders)
		order.Get("/repairer", rbac.PermInterceptor("order.viewfix"), getRepairerOrders)
		order.Get("/repairer/{id:uint}", rbac.PermInterceptor("order.viewall"), forceGetRepairerOrders)
		order.Get("/all", rbac.PermInterceptor("order.viewall"), getAllOrders)
		order.Post("/", rbac.PermInterceptor("order.create"), createOrder)

		order.PartyFunc("/{id:uint}", func(orderID iris.Party) {
			orderID.Get("/", rbac.PermInterceptor("order.viewall"), getOrderByID)
			orderID.Put("/", rbac.PermInterceptor("order.update"), updateOrder)
			orderID.Put("/force", rbac.PermInterceptor("order.updateall"), forceUpdateOrder)
			orderID.Post("/consume", rbac.PermInterceptor("item.consume"), consumeItem)
			// change order status
			orderID.Post("/release", rbac.PermInterceptor("order.update"), releaseOrder)
			orderID.Post("/assign", rbac.PermInterceptor("order.assign"), assignOrder)
			orderID.Post("/selfassign", rbac.PermInterceptor("order.selfassign"), selfAssignOrder)
			orderID.Post("/complete", rbac.PermInterceptor("order.complete"), completeOrder)
			orderID.Post("/cancel", rbac.PermInterceptor("order.cancel"), cancelOrder)
			orderID.Post("/reject", rbac.PermInterceptor("order.reject"), rejectOrder)
			orderID.Post("/report", rbac.PermInterceptor("order.report"), reportOrder)
			orderID.Post("/hold", rbac.PermInterceptor("order.hold"), holdOrder)
			orderID.Post("/appraise", rbac.PermInterceptor("order.appraise"), appraiseOrder)

			orderID.PartyFunc("/comment", func(comment iris.Party) {
				comment.Get("/", rbac.PermInterceptor("comment.view"), getCommentsByOrder)
				comment.Get("/force", rbac.PermInterceptor("comment.viewall"), forceGetCommentsByOrder)
				comment.Post("/", rbac.PermInterceptor("comment.create"), createComment)
				comment.Post("/force", rbac.PermInterceptor("comment.createall"), forceCreateComment)
			})
		})
	})

	mctx.Route.PartyFunc("/tag", func(tag iris.Party) {
		tag.Get("/{id:uint}", middleware.LoginInterceptor, getTagByID)
		tag.Get("/sort", middleware.LoginInterceptor, getAllTagSorts)
		tag.Get("/sort/{name:string}", middleware.LoginInterceptor, getAllTagsBySort)
		tag.Post("/", rbac.PermInterceptor("tag.create"), createTag)
		tag.Delete("/{id:uint}", rbac.PermInterceptor("tag.delete"), deleteTag)
	})

	mctx.Route.PartyFunc("/item", func(item iris.Party) {
		item.Get("/name/{name:string}", rbac.PermInterceptor("item.viewall"), getItemByName)
		item.Get("/name/{name:string}/fuzzy", rbac.PermInterceptor("item.viewall"), getItemsByFuzzyName)
		item.Get("/all", rbac.PermInterceptor("item.viewall"), getAllItems)
		item.Get("/{id:uint}", rbac.PermInterceptor("item.viewall"), getItemByID)
		item.Post("/", rbac.PermInterceptor("item.create"), createItem)
		item.Post("/{id:uint}", rbac.PermInterceptor("item.update"), addItem)
		item.Delete("/{id:uint}", rbac.PermInterceptor("item.delete"), deleteItem)
	})

	mctx.Route.PartyFunc("/comment", func(comment iris.Party) {
		comment.Delete("/{id:uint}", rbac.PermInterceptor("comment.delete"), deleteComment)
		comment.Delete("/{id:uint}/force", rbac.PermInterceptor("comment.deleteall"), forceDeleteComment)
	})
}

// getWxStatusTemplateID godoc
// @Summary 获取 微信 订单状态提醒 模板ID
// @Description 获取 微信 订单状态提醒 模板ID
// @Tags order
// @Produce text/plain
// @Success 200 {string} string "模板ID"
// @Router /wxtmpl/status [get]
func getWxStatusTemplateID(ctx iris.Context) {
	ctx.WriteString(orderConfig.GetString("notify.wechat.status.tmpl"))
}

// getWxCommentTemplateID godoc
// @Summary 获取 微信 订单留言提醒 模板ID
// @Description 获取 微信 订单留言提醒 模板ID
// @Tags order
// @Produce text/plain
// @Success 200 {string} string "模板ID"
// @Router /wxtmpl/comment [get]
func getWxCommentTemplateID(ctx iris.Context) {
	ctx.WriteString(orderConfig.GetString("notify.wechat.comment.tmpl"))
}
