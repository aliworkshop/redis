package redis

import (
	"context"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func (r *repo) Store(ctx context.Context, key string, value any, expiration ...time.Duration) errors.ErrorModel {
	var exp time.Duration = redis.KeepTTL
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	err := r.getTx().Set(ctx, key, value, exp).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) Lock(ctx context.Context, key string, expiration time.Duration) (bool, errors.ErrorModel) {
	ok, err := r.getTx().SetNX(ctx, key, true, expiration).Result()
	if err != nil {
		return false, errors.Internal(err)
	}
	if !ok {
		return false, AlreadyLockedErr
	}
	return ok, nil
}

func (r *repo) Unlock(ctx context.Context, key string) errors.ErrorModel {
	err := r.getTx().Del(ctx, key).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}