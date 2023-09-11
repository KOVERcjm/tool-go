package redis

import (
	"crypto/tls"

	"github.com/go-redis/redis"

	"github.com/kovercjm/tool-go/repository"
)

type Repository struct {
	client *redis.Client
}

func New(config *repository.Config) (*Repository, error) {
	redisOptions := &redis.Options{
		Addr:     config.CacheConfig.Address,
		Password: config.CacheConfig.Password,
		DB:       config.CacheConfig.Database,
	}
	if config.CacheConfig.EnableTLS {
		redisOptions.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}
	client := redis.NewClient(redisOptions)
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Repository{client: client}, nil
}
