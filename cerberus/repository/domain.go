package repository

import (
	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
)

type IDomainRepository interface {
	IRepository[*model.Domain]
}

type DomainRepository struct {
	*Repository[*model.Domain]
}

func NewDomainRepository(pool *database.Pool) *DomainRepository {
	return &DomainRepository{
		Repository: NewRepository[*model.Domain](pool),
	}
}
