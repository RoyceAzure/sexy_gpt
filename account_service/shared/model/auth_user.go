package model

import "github.com/google/uuid"

type AuthUser struct {
	UserId     uuid.UUID
	UserName   string
	Email      string
	IsInternal bool
	RoleName   string
	RoleId     uuid.UUID
	Functions  []string
}
