package handler

import (
	"crowdfunding-minpro-alterra/modules/chat"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatUseCase chat.ChatUseCase
}

func NewChatHandler(chatUC chat.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUC,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	payload, err := chat.ParseChatPayload(c.Request)
	if err != nil {
			c.String(http.StatusBadRequest, "Failed to parse request payload")
			return
	}

	response, err := h.chatUseCase.CompleteChat(payload)
	if err != nil {
			c.String(http.StatusInternalServerError, "Failed to complete chat")
			return
	}

	c.Data(http.StatusOK, "application/json", response)
}
