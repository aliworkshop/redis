package redis

import (
	"context"
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"github.com/redis/go-redis/v9"
)

type repo struct {
	config Config
	client *redis.Client
}

func NewRedisRepository(registry configer.Registry) dbcore.Cache {
	r := new(repo)
	err := registry.Unmarshal(&r.config)
	if err != nil {
		panic(err)
	}

	return r
}

func (r *repo) Initialize() error.ErrorModel {
	client := redis.NewClient(&redis.Options{
		Addr:     r.config.Addr,
		Password: r.config.Password,
		DB:       r.config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	r.client = client
	return nil
}

func (r *repo) GetDB() any {
	return r.client
}

func (r *repo) Ping(ctx context.Context) error.ErrorModel {
	err := r.client.Ping(ctx).Err()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	return nil
}
