package handler

import (
	"ledgerA/internal/dto"
	"ledgerA/internal/middleware"
	"ledgerA/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatHandler handles chatbot endpoints.
type ChatHandler struct {
	chatService service.ChatService
	userService service.UserService
}

// NewChatHandler creates ChatHandler.
func NewChatHandler(chatService service.ChatService, userService service.UserService) *ChatHandler {
	return &ChatHandler{chatService: chatService, userService: userService}
}

// Chat handles POST /api/v1/chat.
func (h *ChatHandler) Chat(c *gin.Context) {
	uid, err := middleware.GetFirebaseUID(c)
	if err != nil {
		dto.Error(c, 401, "ERR_UNAUTHORIZED", "unauthorized")
		return
	}
	user, err := h.userService.GetMe(c.Request.Context(), uid)
	if err != nil {
		dto.Error(c, 404, "ERR_NOT_FOUND", "user not found")
		return
	}

	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, 422, "ERR_VALIDATION", "message is required")
		return
	}

	result, err := h.chatService.Process(c.Request.Context(), user.ID, req.Message)
	if err != nil {
		dto.Error(c, 500, "ERR_INTERNAL", "failed to process message")
		return
	}

	dto.Success(c, result)
}
