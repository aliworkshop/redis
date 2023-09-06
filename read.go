package redis

import (
	"context"
	"encoding/json"
	"github.com/aliworkshop/error"
	"time"
)

func (r *repo) ListKeys(ctx context.Context, pattern string) ([]string, error.ErrorModel) {
	keys := make([]string, 0)
	itr := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	if err := itr.Err(); err != nil {
		return nil, error.DefaultInternalError.WithError(err)
	}
	for itr.Next(ctx) {
		keys = append(keys, itr.Val())
	}
	return keys, nil
}

func (r *repo) Fetch(ctx context.Context, key string) ([]byte, error.ErrorModel) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, error.DefaultInternalError.WithError(err)
	}
	return data, nil
}

func (r *repo) Load(ctx context.Context, key string, result any) error.ErrorModel {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}

	return nil
}

func (r *repo) GetExpiration(ctx context.Context, key string) (time.Duration, error.ErrorModel) {
	dur, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, error.DefaultInternalError.WithError(err)
	}
	return dur, nil
}
