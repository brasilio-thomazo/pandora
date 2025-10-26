package controller

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/service"
	"github.com/gofiber/fiber/v2"
)

type DomainController struct {
	*Controller[*model.Domain]
}

func NewDomainController(db *database.Pool, dbc *database.Cache) *DomainController {
	srv := service.NewDomainService(db, dbc)
	return &DomainController{
		Controller: NewController(srv),
	}
}

func (c *DomainController) Register(router fiber.Router) {
	r := router.Group("/domain")
	r.Get("/", c.index)
	r.Get("/:id", c.show)
	r.Post("/", c.store)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}
