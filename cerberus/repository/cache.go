package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/brasilio/pandora/cerberus/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ICacheRepository[T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindOne(ctx context.Context, id uuid.UUID) (*T, error)
	SetOne(ctx context.Context, id uuid.UUID, data T) error
	SetList(ctx context.Context, data []T) error
	DeleteOne(ctx context.Context, id uuid.UUID) error
	DeleteAll(ctx context.Context) error
}

type CacheRepository[T any] struct {
	rc  *redis.Client
	ttl int64
	key string
}

func NewCacheRepository[T any](cache *database.Cache, key string) *CacheRepository[T] {
	return &CacheRepository[T]{
		rc:  cache.Client,
		ttl: cache.TTL,
		key: key,
	}
}

func (r *CacheRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	key := fmt.Sprintf("%s:all", r.key)
	str, err := r.rc.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var data []T
	err = json.Unmarshal([]byte(str), &data)
	return data, err
}

func (r *CacheRepository[T]) FindOne(ctx context.Context, id uuid.UUID) (*T, error) {
	key := fmt.Sprintf("%s:%s", r.key, id.String())
	str, err := r.rc.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var data T
	err = json.Unmarshal([]byte(str), &data)
	return &data, err
}

func (r *CacheRepository[T]) SetList(ctx context.Context, data []T) error {
	key := fmt.Sprintf("%s:all", r.key)
	str, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.rc.Set(ctx, key, str, time.Second*time.Duration(r.ttl)).Err()
}

func (r *CacheRepository[T]) SetOne(ctx context.Context, id uuid.UUID, data T) error {
	key := fmt.Sprintf("%s:%s", r.key, id.String())
	str, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.rc.Set(ctx, key, str, time.Second*time.Duration(r.ttl)).Err()
}

func (r *CacheRepository[T]) DeleteOne(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("%s:%s", r.key, id.String())
	return r.rc.Del(ctx, key).Err()
}

func (r *CacheRepository[T]) DeleteAll(ctx context.Context) error {
	key := fmt.Sprintf("%s:all", r.key)
	return r.rc.Del(ctx, key).Err()
}
