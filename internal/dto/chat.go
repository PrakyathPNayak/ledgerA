package dto

// ChatRequest contains the user's chat message.
type ChatRequest struct {
	Message string `json:"message" binding:"required" validate:"required,min=1,max=1000"`
}

// ChatAction describes a parsed action from the user's message.
type ChatAction struct {
	Type         string  `json:"type"`
	Name         string  `json:"name,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Account      string  `json:"account,omitempty"`
	Category     string  `json:"category,omitempty"`
	Date         string  `json:"date,omitempty"`
}

// ChatResponse contains the assistant reply and any executed action.
type ChatResponse struct {
	Reply   string      `json:"reply"`
	Action  *ChatAction `json:"action,omitempty"`
	Success bool        `json:"success"`
}
