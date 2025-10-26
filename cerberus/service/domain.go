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

func NewDomainService(pool *database.Pool) *DomainService {
	repo := repository.NewDomainRepository(pool)
	return &DomainService{
		Service: NewService(repo),
		repo:    repo,
	}
}
