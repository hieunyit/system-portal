package dto

// LoginRequest contains credentials for authentication.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse represents returned JWT tokens.
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshRequest contains a refresh token for token renewal.
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
