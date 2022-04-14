package order

import (
	"github.com/xaxys/maintainman/core/middleware"
	"github.com/xaxys/maintainman/module"

	"github.com/kataras/iris/v12"
)

var Module = module.Module{
	ModuleName:    "announce",
	ModuleVersion: "1.0.0",
	ModuleConfig:  orderConfig,
	ModuleEnv: map[string]any{
		"orm.model": []any{
			&Order{},
			&Status{},
			&Tag{},
			&Comment{},
			&Item{},
			&ItemLog{},
		},
	},
	ModuleExport: map[string]any{},
	EntryPoint:   entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	mctx.Scheduler.Every(orderConfig.GetString("appraise.purge")).SingletonMode().Do(autoAppraiseOrderService)
	mctx.Route.PartyFunc("/order", func(order iris.Party) {
		order.Get("/user", middleware.PermInterceptor("order.view"), getUserOrders)
		order.Get("/repairer", middleware.PermInterceptor("order.viewfix"), getRepairerOrders)
		order.Get("/repairer/{id:uint}", middleware.PermInterceptor("order.viewall"), forceGetRepairerOrders)
		order.Get("/all", middleware.PermInterceptor("order.viewall"), getAllOrders)
		order.Post("/", middleware.PermInterceptor("order.create"), createOrder)

		order.PartyFunc("/{id:uint}", func(orderID iris.Party) {
			orderID.Get("/", middleware.PermInterceptor("order.viewall"), getOrderByID)
			orderID.Put("/update", middleware.PermInterceptor("order.update"), updateOrder)
			orderID.Put("/update/force", middleware.PermInterceptor("order.updateall"), forceUpdateOrder)
			orderID.Post("/consume", middleware.PermInterceptor("item.consume"), consumeItem)
			// change order status
			orderID.Post("/release", middleware.PermInterceptor("order.update"), releaseOrder)
			orderID.Post("/assign", middleware.PermInterceptor("order.assign"), assignOrder)
			orderID.Post("/selfassign", middleware.PermInterceptor("order.selfassign"), selfAssignOrder)
			orderID.Post("/complete", middleware.PermInterceptor("order.complete"), completeOrder)
			orderID.Post("/cancel", middleware.PermInterceptor("order.cancel"), cancelOrder)
			orderID.Post("/reject", middleware.PermInterceptor("order.reject"), rejectOrder)
			orderID.Post("/report", middleware.PermInterceptor("order.report"), reportOrder)
			orderID.Post("/hold", middleware.PermInterceptor("order.hold"), holdOrder)
			orderID.Post("/appraise", middleware.PermInterceptor("order.appraise"), appraiseOrder)

			orderID.PartyFunc("/comment", func(comment iris.Party) {
				comment.Get("/", middleware.PermInterceptor("comment.view"), getCommentsByOrder)
				comment.Get("/force", middleware.PermInterceptor("comment.viewall"), forceGetCommentsByOrder)
				comment.Post("/", middleware.PermInterceptor("comment.create"), createComment)
				comment.Post("/force", middleware.PermInterceptor("comment.createall"), forceCreateComment)
			})
		})
	})

	mctx.Route.PartyFunc("/tag", func(tag iris.Party) {
		tag.Get("/{id:uint}", middleware.LoginInterceptor, getTagByID)
		tag.Get("/sort", middleware.LoginInterceptor, getAllTagSorts)
		tag.Get("/sort/{name:string}", middleware.LoginInterceptor, getAllTagsBySort)
		tag.Post("/", middleware.PermInterceptor("tag.create"), createTag)
		tag.Delete("/{id:uint}", middleware.PermInterceptor("tag.delete"), deleteTag)
	})

	mctx.Route.PartyFunc("/item", func(item iris.Party) {
		item.Get("/name/{name:string}", middleware.PermInterceptor("item.viewall"), getItemByName)
		item.Get("/name/{name:string}/fuzzy", middleware.PermInterceptor("item.viewall"), getItemsByFuzzyName)
		item.Get("/all", middleware.PermInterceptor("item.viewall"), getAllItems)
		item.Get("/{id:uint}", middleware.PermInterceptor("item.viewall"), getItemByID)
		item.Post("/", middleware.PermInterceptor("item.create"), createItem)
		item.Post("/{id:uint}", middleware.PermInterceptor("item.update"), addItem)
		item.Delete("/{id:uint}", middleware.PermInterceptor("item.delete"), deleteItem)
	})

	mctx.Route.PartyFunc("/comment", func(comment iris.Party) {
		comment.Delete("/{id:uint}", middleware.PermInterceptor("comment.delete"), deleteComment)
		comment.Delete("/{id:uint}/force", middleware.PermInterceptor("comment.deleteall"), forceDeleteComment)
	})
}
