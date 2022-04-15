package router

import (
	"github.com/xaxys/maintainman/core/controller"
	"github.com/xaxys/maintainman/core/middleware"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

var APIRoute iris.Party

func Register(app *iris.Application) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORS)
	app.AllowMethods(iris.MethodOptions)

	app.PartyFunc("/", func(home iris.Party) {
		home.HandleDir("/", "./assets")
		home.Get("/", func(ctx iris.Context) {
			ctx.Redirect("/index.html")
		})
	})

	v1 := app.Party("/v1")
	v1.Use(middleware.HeaderExtractor, middleware.TokenValidator)
	v1.Done(middleware.ResponseHandler)
	v1.SetExecutionRules(iris.ExecutionRules{Done: iris.ExecutionOptions{Force: true}})
	APIRoute = v1

	v1.Post("/login", middleware.PermInterceptor("user.login"), controller.UserLogin)
	v1.Post("/wxlogin", middleware.PermInterceptor("user.wxlogin"), controller.WxUserLogin)
	v1.Post("/register", middleware.PermInterceptor("user.register"), controller.UserRegister)
	v1.Post("/wxregister", middleware.PermInterceptor("user.wxregister"), controller.WxUserRegister)
	v1.PartyFunc("/", func(api iris.Party) {
		api.Get("/renew", middleware.PermInterceptor("user.renew"), controller.UserRenew)

		api.PartyFunc("/user", func(user iris.Party) {
			user.Get("/", middleware.PermInterceptor("user.view"), controller.GetUser)
			user.Put("/", middleware.PermInterceptor("user.update"), controller.UpdateUser)
			user.Post("/", middleware.PermInterceptor("user.create"), controller.CreateUser)
			user.Get("/all", middleware.PermInterceptor("user.viewall"), controller.GetAllUsers)
			user.Get("/{id:uint}", middleware.PermInterceptor("user.viewall"), controller.GetUserByID)
			user.Put("/{id:uint}", middleware.PermInterceptor("user.updateall"), controller.ForceUpdateUser)
			user.Delete("/{id:uint}", middleware.PermInterceptor("user.delete"), controller.ForceDeleteUser)
			user.Get("/division/{id:uint}", middleware.PermInterceptor("user.viewall"), controller.GetUsersByDivision)
		})

		api.PartyFunc("/role", func(role iris.Party) {
			role.Get("/", middleware.PermInterceptor("role.view"), controller.GetRole)
			role.Post("/", middleware.PermInterceptor("role.create"), controller.CreateRole)
			role.Get("/all", middleware.PermInterceptor("role.viewall"), controller.GetAllRoles)
			role.Get("/{name:string}", middleware.PermInterceptor("role.viewall"), controller.GetRoleByName)
			role.Post("/{name:string}/default", middleware.PermInterceptor("role.update"), controller.SetDefaultRole)
			role.Post("/{name:string}/guest", middleware.PermInterceptor("role.update"), controller.SetGuestRole)
			role.Put("/{name:string}", middleware.PermInterceptor("role.update"), controller.UpdateRole)
			role.Delete("/{name:string}", middleware.PermInterceptor("role.delete"), controller.DeleteRole)
		})

		api.PartyFunc("/permission", func(perm iris.Party) {
			perm.Get("/all", middleware.PermInterceptor("permission.viewall"), controller.GetAllPermissions)
			perm.Get("/{name:string}", middleware.PermInterceptor("permission.viewall"), controller.GetPermission)
		})

		api.PartyFunc("/division", func(division iris.Party) {
			division.Get("/{id:uint}", middleware.PermInterceptor("division.viewall"), controller.GetDivision)
			division.Get("/{id:uint}/children", middleware.PermInterceptor("division.viewall"), controller.GetDivisionsByParentID)
			division.Post("/", middleware.PermInterceptor("division.create"), controller.CreateDivision)
			division.Put("/{id:uint}", middleware.PermInterceptor("division.update"), controller.UpdateDivision)
			division.Delete("/{id:uint}", middleware.PermInterceptor("division.delete"), controller.DeleteDivision)
		})
	})
}
