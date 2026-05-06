package redis

import (
	"context"
	"encoding/json"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func (r *repo) ListKeys(ctx context.Context, pattern string) ([]string, errors.ErrorModel) {
	keys := make([]string, 0)
	itr := r.getLoadTx().Scan(ctx, 0, pattern, 0).Iterator()
	if err := itr.Err(); err != nil {
		return nil, errors.Internal(err)
	}
	for itr.Next(ctx) {
		keys = append(keys, itr.Val())
	}
	return keys, nil
}

func (r *repo) Fetch(ctx context.Context, key string) ([]byte, errors.ErrorModel) {
	data, err := r.getLoadTx().Get(ctx, key).Bytes()
	if err != nil {
		if e, ok := err.(redis.Error); ok && e == redis.Nil {
			return nil, nil
		}
		return nil, errors.Internal(err)
	}
	return data, nil
}

func (r *repo) Load(ctx context.Context, key string, result any) errors.ErrorModel {
	data, err := r.getLoadTx().Get(ctx, key).Bytes()
	if err != nil {
		if e, ok := err.(redis.Error); ok && e == redis.Nil {
			return nil
		}
		return errors.Internal(err)
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return errors.Internal(err)
	}

	return nil
}

func (r *repo) GetExpiration(ctx context.Context, key string) (time.Duration, errors.ErrorModel) {
	dur, err := r.getLoadTx().TTL(ctx, key).Result()
	if err != nil {
		return 0, errors.Internal(err)
	}
	return dur, nil
}