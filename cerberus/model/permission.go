package model

import "github.com/google/uuid"

const (
	Read   = "read"
	Write  = "write"
	Update = "update"
	Delete = "delete"
)

type Permission struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
	DeletedAt   *int64    `json:"deleted_at"`
}
