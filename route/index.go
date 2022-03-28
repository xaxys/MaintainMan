package route

import (
	"maintainman/controller"
	"maintainman/middleware"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// Route ...
func Route(app *iris.Application) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORS)
	app.AllowMethods(iris.MethodOptions)

	app.PartyFunc("/", func(home iris.Party) {
		home.HandleDir("/", "./assets")

		home.Get("/", func(ctx iris.Context) {
			ctx.Redirect("/index.html")
		})

		app.PartyFunc("/v1", func(v1 iris.Party) {
			v1.Use(middleware.HeaderExtractor, middleware.TokenValidator)
			v1.Done(middleware.ResponseHandler)
			v1.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})

			v1.Post("/login", middleware.PermInterceptor("user.login"), controller.UserLogin)
			v1.Post("/wxlogin", middleware.PermInterceptor("user.login"), controller.WxUserLogin)
			v1.Post("/register", middleware.PermInterceptor("user.register"), controller.UserRegister)
			v1.PartyFunc("/", func(api router.Party) {
				api.Get("/renew", middleware.PermInterceptor("user.renew"), controller.UserRenew)

				api.PartyFunc("/user", func(user router.Party) {
					user.Get("/", middleware.PermInterceptor("user.view"), controller.GetUser)
					user.Put("/", middleware.PermInterceptor("user.update"), controller.UpdateUser)
					user.Post("/", middleware.PermInterceptor("user.create"), controller.CreateUser)
					user.Get("/all", middleware.PermInterceptor("user.viewall"), controller.GetAllUsers)
					user.Get("/{id:uint}", middleware.PermInterceptor("user.viewall"), controller.GetUserByID)
					user.Put("/{id:uint}", middleware.PermInterceptor("user.updateall"), controller.ForceUpdateUser)
					user.Delete("/{id:uint}", middleware.PermInterceptor("user.delete"), controller.ForceDeleteUser)
				})

				api.PartyFunc("/role", func(role router.Party) {
					role.Get("/", middleware.PermInterceptor("role.view"), controller.GetRole)
					role.Post("/", middleware.PermInterceptor("role.create"), controller.CreateRole)
					role.Get("/all", middleware.PermInterceptor("role.viewall"), controller.GetAllRoles)
					role.Get("/{name:string}", middleware.PermInterceptor("role.viewall"), controller.GetRoleByName)
					role.Post("/{name:string}/default", middleware.PermInterceptor("role.update"), controller.SetDefaultRole)
					role.Post("/{name:string}/guest", middleware.PermInterceptor("role.update"), controller.SetGuestRole)
					role.Put("/{name:string}", middleware.PermInterceptor("role.update"), controller.UpdateRole)
					role.Delete("/{name:string}", middleware.PermInterceptor("role.delete"), controller.DeleteRole)
				})

				api.PartyFunc("/permission", func(perm router.Party) {
					perm.Get("/all", middleware.PermInterceptor("permission.viewall"), controller.GetAllPermissions)
					perm.Get("/{name:string}", middleware.PermInterceptor("permission.viewall"), controller.GetPermission)
				})

				api.PartyFunc("/announce", func(announce router.Party) {
					announce.Get("/", middleware.PermInterceptor("announce.view"), controller.GetLatestAnnounces)
					announce.Get("/all", middleware.PermInterceptor("announce.viewall"), controller.GetAllAnnounces)
					announce.Get("/{id:uint}", middleware.PermInterceptor("announce.viewall"), controller.GetAnnounce)
					announce.Post("/", middleware.PermInterceptor("announce.create"), controller.CreateAnnounce)
					announce.Put("/{id:uint}", middleware.PermInterceptor("announce.update"), controller.UpdateAnnounce)
					announce.Delete("/{id:uint}", middleware.PermInterceptor("announce.delete"), controller.DeleteAnnounce)
					announce.Get("/{id:uint}/hit", middleware.PermInterceptor("announce.hit"), controller.HitAnnounce)
				})

				api.PartyFunc("/order", func(order router.Party) {
					order.Get("/user", middleware.PermInterceptor("order.view"), controller.GetUserOrders)
					order.Get("/repairer", middleware.PermInterceptor("order.viewfix"), controller.GetRepairerOrders)
					order.Get("/repairer/{id:uint}", middleware.PermInterceptor("order.viewall"), controller.ForceGetRepairerOrders)
					order.Get("/all", middleware.PermInterceptor("order.viewall"), controller.GetAllOrders)
					order.Post("/", middleware.PermInterceptor("order.create"), controller.CreateOrder)

					order.PartyFunc("/{id:uint}", func(orderID router.Party) {
						order.Get("/", middleware.PermInterceptor("order.viewall"), controller.GetOrderByID)
						order.Put("/update", middleware.PermInterceptor("order.update"), controller.UpdateOrder)
						order.Put("/update/force", middleware.PermInterceptor("order.updateall"), controller.ForceUpdateOrder)
						order.Post("/consume", middleware.PermInterceptor("item.consume"), controller.ConsumeItem)
						// change order status
						order.Post("/release", middleware.PermInterceptor("order.update"), controller.ReleaseOrder)
						order.Post("/assign", middleware.PermInterceptor("order.assign"), controller.AssignOrder)
						order.Post("/selfassign", middleware.PermInterceptor("order.selfassign"), controller.SelfAssignOrder)
						order.Post("/complete", middleware.PermInterceptor("order.complete"), controller.CompleteOrder)
						order.Post("/cancel", middleware.PermInterceptor("order.cancel"), controller.CancelOrder)
						order.Post("/reject", middleware.PermInterceptor("order.reject"), controller.RejectOrder)
						order.Post("/report", middleware.PermInterceptor("order.report"), controller.ReportOrder)
						order.Post("/hold", middleware.PermInterceptor("order.hold"), controller.HoldOrder)
						order.Post("/appraise", middleware.PermInterceptor("order.appraise"), controller.AppraiseOrder)

						order.PartyFunc("/comment", func(comment router.Party) {
							comment.Get("/", middleware.PermInterceptor("comment.view"), controller.GetCommentsByOrder)
							comment.Get("/force", middleware.PermInterceptor("comment.viewall"), controller.ForceGetCommentsByOrder)
							comment.Post("/", middleware.PermInterceptor("comment.create"), controller.CreateComment)
							comment.Post("/force", middleware.PermInterceptor("comment.createall"), controller.ForceCreateComment)
						})
					})
				})

				api.PartyFunc("/tag", func(tag router.Party) {
					tag.Get("/{id:uint}", middleware.PermInterceptor("tag.viewall"), controller.GetTagByID)
					tag.Get("/sort", middleware.PermInterceptor("tag.viewall"), controller.GetAllTagSorts)
					tag.Get("/sort/{name:string}", middleware.PermInterceptor("tag.viewall"), controller.GetAllTagsBySort)
					tag.Post("/", middleware.PermInterceptor("tag.create"), controller.CreateTag)
					tag.Delete("/{id:uint}", middleware.PermInterceptor("tag.delete"), controller.DeleteTag)
				})

				api.PartyFunc("/item", func(item router.Party) {
					item.Get("/name/{name:string}", middleware.PermInterceptor("item.viewall"), controller.GetItemByName)
					item.Get("/name/{name:string}/fuzzy", middleware.PermInterceptor("item.viewall"), controller.GetItemsByFuzzyName)
					item.Get("/all", middleware.PermInterceptor("item.viewall"), controller.GetAllItems)
					item.Get("/{id:uint}", middleware.PermInterceptor("item.viewall"), controller.GetItemByID)
					item.Post("/", middleware.PermInterceptor("item.create"), controller.CreateItem)
					item.Post("/{id:uint}", middleware.PermInterceptor("item.update"), controller.AddItem)
					item.Delete("/{id:uint}", middleware.PermInterceptor("item.delete"), controller.DeleteItem)
				})

				api.PartyFunc("/comment", func(comment router.Party) {
					comment.Delete("/{id:uint}", middleware.PermInterceptor("comment.delete"), controller.DeleteComment)
					comment.Delete("/{id:uint}/force", middleware.PermInterceptor("comment.deleteall"), controller.ForceDeleteComment)
				})
			})
		})
	})
}
