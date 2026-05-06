package redis

import (
	"context"
	"github.com/aliworkshop/errors"
)

func (r *repo) Exists(ctx context.Context, key string) (bool, errors.ErrorModel) {
	val, err := r.getTx().Exists(ctx, key).Result()
	if err != nil {
		return false, errors.Internal(err)
	}

	return val > 0, nil
}