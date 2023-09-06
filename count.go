package redis

import (
	"context"
	"github.com/aliworkshop/error"
)

func (r *repo) Count(ctx context.Context, pattern string) (uint64, error.ErrorModel) {
	var cursor uint64
	err := r.client.Scan(ctx, cursor, pattern, 0).Err()
	if err != nil {
		return 0, error.DefaultInternalError.WithError(err)
	}
	return cursor, nil
}
