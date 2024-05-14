package app

import (
	"roby-backend-golang/api"
	"roby-backend-golang/app/modules"
	"roby-backend-golang/config"
	"roby-backend-golang/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "roby-backend-golang/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Run(config *config.AppConfig, dbCon *utils.DatabaseConnection) (*fiber.App, string) {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] [${ip}:${port}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2 Jan 2006 15:04:05",
	}))
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API IS UP!!!!",
		})
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	controller := modules.RegistrationModules(dbCon, config)
	api.RegistrationPath(app, controller)
	return app, config.App.Port
}
