package server

import (
	"fmt"
	"log"

	"bitbucket.org/brasilio/pandora/cerberus/config"
	"bitbucket.org/brasilio/pandora/cerberus/controller"
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"github.com/gofiber/fiber/v2"
)

func Httpd(config *config.Config, db *database.Pool) error {
	app := fiber.New(fiber.Config{Prefork: false})
	addr := fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort)
	log.Printf("listening on %s", addr)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "service is healthy"})
	})

	domain := controller.NewDomainController(db)
	domain.Register(app.Group("/api"))

	return app.Listen(addr)
}
