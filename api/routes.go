package api

import (
	"roby-backend-golang/api/middlewares"
	"roby-backend-golang/api/user"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	UserController *user.Controller
}

func RegistrationPath(e *fiber.App, controller Controller) {
	route := e.Group("/v1")
	routeUser := route.Group("/user")
	routeUser.Post("/login", controller.UserController.Login)
	routeUser.Post("/register", controller.UserController.Register)

	routeUser.Use(middlewares.MiddleJWT)
	routeUser.Delete("/logout", controller.UserController.Logout)
	routeUser.Get("/me", controller.UserController.GetMe)
	routeUser.Get("/find-random", controller.UserController.GetRandomUser)
	routeUser.Post("/swipe", controller.UserController.SwipeUser)

	routePackage := route.Group("/package")
	routePackage.Get("/list", controller.UserController.GetListPackage)
	routePackage.Post("/purchase", middlewares.MiddleJWT, controller.UserController.PurchasePackage)
}
