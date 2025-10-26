package service

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/repository"
)

type IPermissionService interface {
	IService[*model.Permission]
}

type PermissionService struct {
	*Service[*model.Permission]
}

func NewPermissionService(db *database.Pool, dbc *database.Cache) *PermissionService {
	repo := repository.NewPermissionRepository(db)
	cache := repository.NewCacheRepository[*model.Permission](dbc, "permission")
	return &PermissionService{
		Service: NewService(repo, cache),
	}
}
