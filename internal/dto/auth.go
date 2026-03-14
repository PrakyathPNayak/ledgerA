package dto

import "ledgerA/internal/model"

// SyncRequest contains the auth sync payload from frontend.
type SyncRequest struct {
	FirebaseToken string  `json:"firebase_token" binding:"required" validate:"required,min=10"`
	DisplayName   string  `json:"display_name" binding:"required" validate:"required,min=1,max=100"`
	Email         string  `json:"email" binding:"required" validate:"required,email,max=255"`
	CurrencyCode  *string `json:"currency_code,omitempty" validate:"omitempty,len=3"`
}

// SyncResponse contains the synchronized user profile.
type SyncResponse struct {
	User UserResponse `json:"user"`
}

// ToModel converts a SyncRequest to a model.User.
func (req SyncRequest) ToModel(firebaseUID string, existingCurrencyCode string) model.User {
	currency := existingCurrencyCode
	if currency == "" && req.CurrencyCode != nil {
		currency = *req.CurrencyCode
	}

	return model.User{
		FirebaseUID:  firebaseUID,
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		CurrencyCode: currency,
	}
}
