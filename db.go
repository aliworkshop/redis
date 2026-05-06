package redis

import (
	"context"
	"crypto/tls"
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
)

type repo struct {
	config Config
	client *redis.Client
	pipe   redis.Pipeliner
	tx     *redis.Tx
}

func NewRedisRepository(registry configer.Registry) dbcore.Cache {
	r := new(repo)
	err := registry.Unmarshal(&r.config)
	if err != nil {
		panic(err)
	}
	r.config.Initialize()

	return r
}

func (r *repo) Initialize() errors.ErrorModel {
	opts := &redis.Options{
		Addr:     r.config.Addr,
		Username: r.config.Username,
		Password: r.config.Password,
		DB:       r.config.DB,
	}
	if r.config.Tls {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: r.config.Tls}
	}
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return errors.Internal(err)
	}
	r.client = client
	return nil
}

func (r *repo) GetDB() any {
	return r.client
}

func (r *repo) Ping(ctx context.Context) errors.ErrorModel {
	err := r.getTx().Ping(ctx).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) getTx() redis.Cmdable {
	if r.pipe != nil {
		return r.pipe
	} else if r.tx != nil {
		return r.tx
	}
	return r.client
}

func (r *repo) getLoadTx() redis.Cmdable {
	if r.tx != nil {
		return r.tx
	}
	return r.client
}