package handlers

import (
	"lite_ai_hub/ai_hub/internal/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandleMessage(c fiber.Ctx) {
	var session dto.Session
	if err := c.Bind().Body(&session); err != nil {
		log.Error(err)
		return
	}
	// do smth with session
}
