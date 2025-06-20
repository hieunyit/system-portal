package dto

import "github.com/google/uuid"

type AuditResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	Action   string    `json:"action"`
	Resource string    `json:"resource"`
	Success  bool      `json:"success"`
}
