package service

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/repository"
)

type IDomainService interface {
	IService[*model.Domain]
}

type DomainService struct {
	*Service[*model.Domain]
	repo repository.IDomainRepository
}

func NewDomainService(db *database.Pool, dbc *database.Cache) *DomainService {
	repo := repository.NewDomainRepository(db)
	cache := repository.NewCacheRepository[*model.Domain](dbc, "domain")
	return &DomainService{
		Service: NewService(repo, cache),
		repo:    repo,
	}
}
