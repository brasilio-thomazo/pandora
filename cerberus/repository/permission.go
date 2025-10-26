package repository

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
)

type PermissionRepository struct {
	*Repository[*model.Permission]
}

func NewPermissionRepository(pool *database.Pool) *PermissionRepository {
	return &PermissionRepository{
		Repository: NewRepository[*model.Permission](pool),
	}
}
