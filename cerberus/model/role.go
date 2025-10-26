package model

import "github.com/google/uuid"

type Role struct {
	Model
	PermissionId uuid.UUID  `json:"permission_id"`
	DomainId     uuid.UUID  `json:"domain_id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description"`
	Permission   Permission `json:"permission"`
	Domain       Domain     `json:"domain"`
}
