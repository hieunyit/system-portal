package dto

import "github.com/google/uuid"

type PortalUserRequest struct {
	Username string    `json:"username" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	FullName string    `json:"fullName"`
	Password string    `json:"password"`
	GroupID  uuid.UUID `json:"groupId"`
}

type PortalUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
	GroupID  uuid.UUID `json:"groupId"`
	IsActive bool      `json:"isActive"`
}
