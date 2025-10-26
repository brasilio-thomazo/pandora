package service

import (
	"context"

	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/repository"
	"github.com/google/uuid"
)

type IRoleService interface {
	IService[*model.Role]
}

type RoleService struct {
	*Service[*model.Role]
}

func NewRoleService(db *database.Pool, dbc *database.Cache) *RoleService {
	repo := repository.NewRoleRepository(db)
	cache := repository.NewCacheRepository[*model.Role](dbc, "role")
	return &RoleService{
		Service: NewService(repo, cache),
	}
}

func (s *RoleService) FindAll(ctx context.Context) ([]*model.Role, error) {
	return s.repo.FindAllWithRelationships(ctx, "Permission", "Domain")
}

func (s *RoleService) FindById(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	return s.repo.FindByIdWithRelationships(ctx, id, "Permission", "Domain")
}
