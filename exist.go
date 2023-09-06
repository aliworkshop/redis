package redis

import (
	"context"
	"github.com/aliworkshop/error"
)

func (r *repo) Exists(ctx context.Context, key string) (bool, error.ErrorModel) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, error.DefaultInternalError.WithError(err)
	}

	return val > 0, nil
}
