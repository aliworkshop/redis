package redis

import (
	"context"
	"github.com/aliworkshop/errors"
)

func (r *repo) Count(ctx context.Context, pattern string) (uint64, errors.ErrorModel) {
	var cursor uint64
	err := r.getTx().Scan(ctx, cursor, pattern, 0).Err()
	if err != nil {
		return 0, errors.Internal(err)
	}
	return cursor, nil
}