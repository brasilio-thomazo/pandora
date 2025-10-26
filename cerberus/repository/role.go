package repository

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
)

type RoleRepository struct {
	*Repository[*model.Role]
}

func NewRoleRepository(pool *database.Pool) *RoleRepository {
	return &RoleRepository{
		Repository: NewRepository[*model.Role](pool),
	}
}
