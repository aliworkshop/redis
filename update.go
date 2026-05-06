package redis

import (
	"context"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
)

func (r *repo) GetInt(ctx context.Context, key string) (int64, errors.ErrorModel) {
	value, err := r.getLoadTx().Get(ctx, key).Int64()
	if err != nil {
		if e, ok := err.(redis.Error); ok && e == redis.Nil {
			return 0, nil
		}
		return 0, errors.Internal(err)
	}
	return value, nil
}

func (r *repo) DecrBy(ctx context.Context, key string, decrement int64) errors.ErrorModel {
	err := r.getTx().DecrBy(ctx, key, decrement).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) IncrBy(ctx context.Context, key string, increment int64) errors.ErrorModel {
	err := r.getTx().IncrBy(ctx, key, increment).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}