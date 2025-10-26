package server

import (
	"fmt"
	"log"

	"bitbucket.org/brasilio/pandora/cerberus/config"
	"bitbucket.org/brasilio/pandora/cerberus/controller"
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"github.com/gofiber/fiber/v2"
)

func Httpd(config *config.Config, db *database.Pool, dbc *database.Cache) error {
	app := fiber.New(fiber.Config{Prefork: false})
	addr := fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort)
	log.Printf("listening on %s", addr)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "service is healthy"})
	})

	domain := controller.NewDomainController(db, dbc)
	domain.Register(app.Group("/api"))

	permission := controller.NewPermissionController(db, dbc)
	permission.Register(app.Group("/api"))

	role := controller.NewRoleController(db, dbc)
	role.Register(app.Group("/api"))

	defer app.Shutdown()

	return app.Listen(addr)
}
