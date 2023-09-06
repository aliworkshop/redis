package redis

import (
	"context"
	"github.com/aliworkshop/error"
	"time"
)

func (r *repo) Expire(ctx context.Context, key string, expiration time.Duration) error.ErrorModel {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, key string) error.ErrorModel {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return error.DefaultInternalError.WithError(err)
	}
	return nil
}
