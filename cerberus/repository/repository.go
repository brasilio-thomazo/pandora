package repository

import (
	"context"
	"time"

	"bitbucket.org/brasilio/pandora/cerberus/database"
	"bitbucket.org/brasilio/pandora/cerberus/model"
	"github.com/google/uuid"
)

type IRepository[T model.TypeModel] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]T, error)
	FindById(ctx context.Context, id uuid.UUID) (T, error)
	ExistsById(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, data T) error
	Update(ctx context.Context, id uuid.UUID, data T) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repository[T model.TypeModel] struct {
	db *database.Pool
}

func NewRepository[T model.TypeModel](pool *database.Pool) *Repository[T] {
	return &Repository[T]{
		db: pool,
	}
}

func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var result []T
	err := r.db.Read.WithContext(ctx).Find(&result).Error
	return result, err
}

func (r *Repository[T]) FindByIds(ctx context.Context, ids []uuid.UUID) ([]T, error) {
	var result []T
	err := r.db.Read.WithContext(ctx).Where("id in (?)", ids).Find(&result).Error
	return result, err
}

func (r *Repository[T]) FindById(ctx context.Context, id uuid.UUID) (T, error) {
	var result T
	err := r.db.Read.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return result, err
}

func (r *Repository[T]) ExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Read.WithContext(ctx).
		Model(new(T)).
		Select("count (*) > 0").
		Where("id = ? and deleted_at is null", id).
		Find(&exists).Error
	return exists, err
}

func (r *Repository[T]) Create(ctx context.Context, data T) error {
	return r.db.Write.WithContext(ctx).Create(&data).Error
}

func (r *Repository[T]) Update(ctx context.Context, id uuid.UUID, data T) error {
	return r.db.Write.WithContext(ctx).Where("id = ?", id).Updates(&data).Error
}

func (r *Repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	timestamp := time.Now().Unix()
	return r.db.Write.WithContext(ctx).
		Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", &timestamp).Error
}
