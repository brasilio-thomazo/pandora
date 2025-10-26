package controller

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/service"
	"github.com/gofiber/fiber/v2"
)

type PermissionController struct {
	*Controller[*model.Permission]
}

func NewPermissionController(db *database.Pool, dbc *database.Cache) *PermissionController {
	srv := service.NewPermissionService(db, dbc)
	return &PermissionController{
		Controller: NewController(srv),
	}
}

func (c *PermissionController) Register(router fiber.Router) {
	r := router.Group("/permission")
	r.Get("/", c.index)
	r.Get("/:id", c.show)
	r.Post("/", c.store)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}
