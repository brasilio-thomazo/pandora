package model

type Domain struct {
	Model
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Paths       []string `json:"paths" gorm:"serializer:json"`
}
