package redis

import (
	"context"
	"github.com/aliworkshop/error"
	"github.com/redis/go-redis/v9"
	"time"
)

func (r *repo) Store(ctx context.Context, key string, value any, expiration ...time.Duration) error.ErrorModel {
	var exp time.Duration = redis.KeepTTL
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	err := r.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	return nil
}

func (r *repo) Lock(ctx context.Context, key string, expiration time.Duration) error.ErrorModel {
	err := r.client.SetNX(ctx, key, true, expiration).Err()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	return nil
}
