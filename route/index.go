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
			v1.Done(middleware.ResponseHandler)
			v1.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})

			v1.Post("/login", middleware.PermInterceptor("user.login"), controller.UserLogin)
			v1.Post("/wxlogin", middleware.PermInterceptor("user.login"), controller.WxUserLogin)
			v1.Post("/register", middleware.PermInterceptor("user.register"), controller.UserRegister)
			v1.PartyFunc("/", func(account router.Party) {
				account.Use(middleware.HeaderExtractor, middleware.TokenValidator, middleware.LoginInterceptor)

				account.Get("/renew", middleware.PermInterceptor("user.renew"), controller.UserRenew)
				account.Get("/graphql", controller.GetGraphQL)

				account.PartyFunc("/user", func(user router.Party) {
					user.Get("/", controller.GetUser)
					user.Put("/", middleware.PermInterceptor("user.update"), controller.UpdateUser)
					user.Post("/", middleware.PermInterceptor("user.create"), controller.CreateUser)
					user.Get("/all", middleware.PermInterceptor("user.viewall"), controller.GetAllUsers)
					user.Get("/{id:uint}", middleware.PermInterceptor("user.viewall"), controller.GetUserByID)
					user.Put("/{id:uint}", middleware.PermInterceptor("user.updateall"), controller.UpdateUserByID)
					user.Delete("/{id:uint}", middleware.PermInterceptor("user.delete"), controller.DeleteUserByID)
				})

				account.PartyFunc("/role", func(role router.Party) {
					role.Get("/", controller.GetRole)
					role.Post("/", middleware.PermInterceptor("role.create"), controller.CreateRole)
					role.Get("/all", middleware.PermInterceptor("role.viewall"), controller.GetAllRoles)
					role.Get("/{name:string}", middleware.PermInterceptor("role.viewall"), controller.GetRoleByName)
					role.Post("/{name:string}/default", middleware.PermInterceptor("role.update"), controller.SetDefaultRole)
					role.Post("/{name:string}/guest", middleware.PermInterceptor("role.update"), controller.SetGuestRole)
					role.Put("/{name:string}", middleware.PermInterceptor("role.update"), controller.UpdateRoleByName)
					role.Delete("/{name:string}", middleware.PermInterceptor("role.delete"), controller.DeleteRoleByName)
				})

				account.PartyFunc("/permission", func(perm router.Party) {
					perm.Get("/all", middleware.PermInterceptor("permission.viewall"), controller.GetAllPermissions)
					perm.Get("/{name:string}", middleware.PermInterceptor("permission.viewall"), controller.GetPermissionByName)
				})

				account.PartyFunc("/announce", func(announce router.Party) {
					announce.Get("/latest/{offset:uint}", middleware.PermInterceptor("announce.viewall"), controller.GetLatestAnnounces)
					announce.Get("/all", middleware.PermInterceptor("announce.viewall"), controller.GetAllAnnounces)
					announce.Get("/{id:uint}", middleware.PermInterceptor("announce.viewall"), controller.GetAnnounceByID)
					announce.Post("/", middleware.PermInterceptor("announce.create"), controller.CreateAnnounceByID)
					announce.Put("/{id:uint}", middleware.PermInterceptor("announce.update"), controller.UpdateAnnounceByID)
					announce.Delete("/{id:uint}", middleware.PermInterceptor("announce.delete"), controller.DeleteAnnounceByID)
					announce.Get("/{id:uint}/hit", middleware.PermInterceptor("announce.viewall"), controller.HitAnnounceByID)
				})

				account.PartyFunc("/order", func(order router.Party) {
					order.PartyFunc("/user", func(user router.Party) {
						order.Get("/{status:uint}/{offset:uint}", controller.GetUserOrders)
						order.Put("/{id:uint}", middleware.PermInterceptor("order.update"), controller.UpdateOrder)
					})

					order.PartyFunc("/repairer", func(repairer router.Party) {
						order.Get("/{current:bool}/{offset:uint}", controller.GetRepairerOrders)
					})

					order.Get("/all", middleware.PermInterceptor("order.viewall"), controller.GetAllOrders)
					order.Get("/{id:uint}", middleware.PermInterceptor("order.viewall"), controller.GetOrderByID)
					order.Post("/", middleware.PermInterceptor("order.create"), controller.CreateOrder)
					order.Put("/{id:uint}", middleware.PermInterceptor("order.updateall"), controller.UpdateOrderByID)
					// change order status
					order.Post("/{id:uint}/release", middleware.PermInterceptor("order.update"), controller.ReleaseOrder)
					order.Post("/{id:uint}/assign/{repairer:uint}", middleware.PermInterceptor("order.assign"), controller.AssignOrder)
					order.Post("/{id:uint}/selfassign", middleware.PermInterceptor("order.selfassign"), controller.SelfAssignOrder)
					order.Post("/{id:uint}/complete", middleware.PermInterceptor("order.complete"), controller.CompleteOrder)
					order.Post("/{id:uint}/cancel", middleware.PermInterceptor("order.cancel"), controller.CancelOrder)
					order.Post("/{id:uint}/reject", middleware.PermInterceptor("order.reject"), controller.RejectOrder)
					order.Post("/{id:uint}/report", middleware.PermInterceptor("order.report"), controller.ReportOrder)
					order.Post("/{id:uint}/hold", middleware.PermInterceptor("order.hold"), controller.HoldOrder)
					order.Post("/{id:uint}/appraise/{appraisal:uint}", middleware.PermInterceptor("order.appraise"), controller.AppraiseOrder)
				})
			})
		})
	})
}
