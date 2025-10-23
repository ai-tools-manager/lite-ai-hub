package handlers

import (
	"lite_ai_hub/ai_hub/internal/dto"

	"github.com/gofiber/fiber/v3"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandleMessage(c *fiber.Ctx) {
	var session dto.Session
	if err := c.BodyParser(&session); err != nil {
		return err
	}
	// do smth with session
}
