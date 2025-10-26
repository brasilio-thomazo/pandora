package database

import (
	"log"

	"bitbucket.org/brasilio/pandora/cerberus/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	TTL    int64
}

func NewCache(cfg *config.Config) *Cache {
	return &Cache{
		Client: get_client(cfg),
		TTL:    cfg.CacheTTL,
	}
}

func (c *Cache) Close() {
	c.Client.Close()
}

func get_client(config *config.Config) *redis.Client {
	log.Printf("connecting to redis at %s", config.RedisAddr)
	cfg := &redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisDatabase,
		Password: config.RedisPassword,
	}
	return redis.NewClient(cfg)
}
