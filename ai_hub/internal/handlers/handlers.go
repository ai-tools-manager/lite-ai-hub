// Package handlers provides HTTP request handlers for the AI Hub API.
package handlers

import (
	"fmt"
	"lite_ai_hub/ai_hub/internal/dto"
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Handlers contains all HTTP request handlers.
type Handlers struct{}

// NewHandlers creates a new Handlers instance.
func NewHandlers() *Handlers {
	return &Handlers{}
}

// CreateChat handles POST /chat requests to create a new chat.
func (h *Handlers) CreateChat(c fiber.Ctx) error {
	chatID := uuid.New()
	return c.JSON(dto.ChatResponse{ChatID: chatID})
}

// GetChatList handles GET /chat/list requests to retrieve all chats.
func (h *Handlers) GetChatList(c fiber.Ctx) error {
	return c.JSON(dto.ChatListResponse{Chats: []dto.ChatResponse{}})
}

// GetMessages handles GET /message/{chat_id} requests to retrieve chat messages.
func (h *Handlers) GetMessages(c fiber.Ctx) error {
	chatID := c.Params("chat_id")
	if _, err := uuid.Parse(chatID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid chat_id"})
	}
	return c.JSON(dto.MessageListResponse{Messages: []dto.Message{}})
}

// SendMessage handles POST /message/{chat_id} requests to send a message to a chat.
func (h *Handlers) SendMessage(c fiber.Ctx) error {
	chatID := c.Params("chat_id")
	if _, err := uuid.Parse(chatID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid chat_id"})
	}

	var req dto.MessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	return c.SendStatus(200)
}

// GetLibraries handles GET /lib/list requests to retrieve all installed libraries.
func (h *Handlers) GetLibraries(c fiber.Ctx) error {
	return c.JSON(dto.LibraryListResponse{})
}

// InstallLibrary handles POST /lib requests to install a new library.
func (h *Handlers) InstallLibrary(c fiber.Ctx) error {
	var req dto.LibraryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	return c.SendStatus(200)
}

// DeleteLibrary handles DELETE /lib/{lib_url} requests to remove a library.
func (h *Handlers) DeleteLibrary(c fiber.Ctx) error {
	libURL, err := url.Parse(c.Params("lib_url"))
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "incorrect lib_url"})
	}
	fmt.Println(libURL)
	return c.SendStatus(200)
}
