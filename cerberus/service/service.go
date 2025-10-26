package service

import (
	"context"
	"time"

	"bitbucket.org/brasilio/pandora/cerberus/model"
	"bitbucket.org/brasilio/pandora/cerberus/repository"
	"github.com/google/uuid"
)

type IService[T model.TypeModel] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]T, error)
	FindById(ctx context.Context, id uuid.UUID) (T, error)
	ExistsById(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, data T) error
	Update(ctx context.Context, id uuid.UUID, data T) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service[T model.TypeModel] struct {
	repo repository.IRepository[T]
}

func NewService[T model.TypeModel](repo repository.IRepository[T]) *Service[T] {
	return &Service[T]{
		repo: repo,
	}
}

func (s *Service[T]) FindAll(ctx context.Context) ([]T, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service[T]) FindByIds(ctx context.Context, ids []uuid.UUID) ([]T, error) {
	return s.repo.FindByIds(ctx, ids)
}

func (s *Service[T]) FindById(ctx context.Context, id uuid.UUID) (T, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service[T]) ExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.ExistsById(ctx, id)
}

func (s *Service[T]) Create(ctx context.Context, data T) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	data.SetId(id)
	timestamp := time.Now().Unix()
	data.SetCreatedAt(timestamp)
	data.SetUpdatedAt(timestamp)
	return s.repo.Create(ctx, data)
}

func (s *Service[T]) Update(ctx context.Context, id uuid.UUID, data T) error {
	timestamp := time.Now().Unix()
	data.SetId(id)
	data.SetUpdatedAt(timestamp)
	return s.repo.Update(ctx, id, data)
}

func (s *Service[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
