package controller

import (
	"fmt"
	"time"

	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Error struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
}

func NewError(message string, code int) *Error {
	return &Error{
		Message:   message,
		Code:      code,
		Timestamp: time.Now().Unix(),
	}
}

func NewFromError(err error) *Error {
	return &Error{
		Message:   err.Error(),
		Code:      500,
		Timestamp: time.Now().Unix(),
	}
}

func BadRequest(err error) *Error {
	return &Error{
		Message:   err.Error(),
		Code:      400,
		Timestamp: time.Now().Unix(),
	}
}

func NotFound(err error) *Error {
	return &Error{
		Message:   err.Error(),
		Code:      404,
		Timestamp: time.Now().Unix(),
	}
}

type IController[T model.TypeModel] interface {
	index(ctx *fiber.Ctx) error
	show(ctx *fiber.Ctx) error
	store(ctx *fiber.Ctx) error
	update(ctx *fiber.Ctx) error
	delete(ctx *fiber.Ctx) error
}

type Controller[T model.TypeModel] struct {
	srv service.IService[T]
}

func NewController[T model.TypeModel](srv service.IService[T]) *Controller[T] {
	return &Controller[T]{
		srv: srv,
	}
}

func (c *Controller[T]) index(ctx *fiber.Ctx) error {
	data, err := c.srv.FindAll(ctx.Context())
	if err != nil {
		return ctx.Status(500).JSON(NewFromError(err))
	}
	return ctx.JSON(data)
}

func (c *Controller[T]) show(ctx *fiber.Ctx) error {
	var id uuid.UUID
	if err := ctx.ParamsParser(&id); err != nil {
		return ctx.Status(400).JSON(BadRequest(err))
	}
	data, err := c.srv.FindById(ctx.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(404).JSON(NewError("not found", 404))
		}
		return ctx.Status(500).JSON(NewFromError(err))
	}
	return ctx.JSON(data)
}

func (c *Controller[T]) store(ctx *fiber.Ctx) error {
	var data T
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(400).JSON(BadRequest(err))

	}
	if err := c.srv.Create(ctx.Context(), data); err != nil {
		return ctx.Status(500).JSON(NewFromError(err))

	}
	uri := fmt.Sprintf("%s/%s", ctx.Path(), data.GetId().String())
	ctx.Set("Location", uri)
	return ctx.Status(201).JSON(data)
}

func (c *Controller[T]) update(ctx *fiber.Ctx) error {
	var id uuid.UUID
	if err := ctx.ParamsParser(&id); err != nil {
		return ctx.Status(400).JSON(BadRequest(err))
	}
	var data T
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(400).JSON(BadRequest(err))
	}
	data.SetId(id)
	if err := c.srv.Update(ctx.Context(), id, data); err != nil {
		return ctx.Status(500).JSON(NewFromError(err))
	}
	return ctx.Status(200).JSON(data)
}

func (c *Controller[T]) delete(ctx *fiber.Ctx) error {
	var id uuid.UUID
	if err := ctx.ParamsParser(&id); err != nil {
		return ctx.Status(400).JSON(BadRequest(err))
	}
	if err := c.srv.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(500).JSON(NewFromError(err))
	}
	return ctx.Status(204).JSON(nil)
}
