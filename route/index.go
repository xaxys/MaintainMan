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
					role.Post("/default/{name:string}", middleware.PermInterceptor("role.update"), controller.SetDefaultRole)
					role.Post("/default/{name:string}", middleware.PermInterceptor("role.update"), controller.SetGuestRole)
					role.Put("/{name:string}", middleware.PermInterceptor("role.update"), controller.UpdateRoleByName)
					role.Delete("/{name:string}", middleware.PermInterceptor("role.delete"), controller.DeleteRoleByName)
				})

				account.PartyFunc("/permission", func(perm router.Party) {
					perm.Get("/all", middleware.PermInterceptor("permission.viewall"), controller.GetAllPermissions)
					perm.Get("/{name:string}", middleware.PermInterceptor("permission.viewall"), controller.GetPermissionByName)
				})
			})
		})
	})
}
