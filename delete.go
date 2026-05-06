package redis

import (
	"context"
	"github.com/aliworkshop/errors"
	"time"
)

func (r *repo) Expire(ctx context.Context, key string, expiration time.Duration) errors.ErrorModel {
	err := r.getTx().Expire(ctx, key, expiration).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, key string) errors.ErrorModel {
	err := r.getTx().Del(ctx, key).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) DeleteWithPattern(ctx context.Context, pattern string) errors.ErrorModel {
	iter := r.getTx().Scan(ctx, 0, pattern, 10).Iterator()
	pipe := r.getTx().Pipeline()
	delCount := 0

	for iter.Next(ctx) {
		key := iter.Val()
		pipe.Del(ctx, key)
		delCount++
	}

	if err := iter.Err(); err != nil {
		return errors.Internal(err)
	}

	if delCount > 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			return errors.Internal(err)
		}
	}

	return nil
}