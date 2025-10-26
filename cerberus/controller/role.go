package controller

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/service"
	"github.com/gofiber/fiber/v2"
)

type RoleController struct {
	*Controller[*model.Role]
}

func NewRoleController(db *database.Pool, dbc *database.Cache) *RoleController {
	srv := service.NewRoleService(db, dbc)
	return &RoleController{
		Controller: NewController(srv),
	}
}

func (c *RoleController) Register(router fiber.Router) {
	r := router.Group("/role")
	r.Get("/", c.index)
	r.Get("/:id", c.show)
	r.Post("/", c.store)
	r.Put("/:id", c.update)
	r.Delete("/:id", c.delete)
}
