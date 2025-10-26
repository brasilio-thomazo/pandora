package model

const (
	Read   = "read"
	Write  = "write"
	Update = "update"
	Delete = "delete"
)

type Permission struct {
	Model
	Name        string   `json:"name"`
	Permissions []string `json:"permissions" gorm:"serializer:json"`
}
