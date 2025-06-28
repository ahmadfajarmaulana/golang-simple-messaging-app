package bootstrap

import (
	"simple-messaging-app/app/ws"
	"simple-messaging-app/pkg/database"
	"simple-messaging-app/pkg/env"
	"simple-messaging-app/pkg/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	database.SetupDatabase()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	go ws.ServeWSMessaging(app)
	router.InstallRouter(app)

	return app
}
