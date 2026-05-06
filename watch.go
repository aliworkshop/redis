package redis

import (
	"context"
	"fmt"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
)

func (r *repo) Watch(ctx context.Context, key string, fn func(cache dbcore.Cache) errors.ErrorModel) errors.ErrorModel {
	for i := 0; i < 5; i++ {
		err := r.client.Watch(ctx, func(tx *redis.Tx) error {
			wrapper := &repo{client: r.client, tx: tx, pipe: tx.TxPipeline()}
			e := fn(wrapper)
			if e != nil {
				return e
			}

			return wrapper.Exec(ctx)
		}, key)

		if err != nil && errors.Is(err, redis.TxFailedErr) {
			// Retry if key was modified during WATCH
			continue
		}
		return errors.HandleError(err)
	}
	return errors.Internal(fmt.Errorf("transaction failed after %d retries", 5))
}

func (r *repo) Exec(ctx context.Context) errors.ErrorModel {
	if r.pipe == nil {
		return errors.Internal().WithMessage("pipeline is nil")
	}
	_, err := r.pipe.Exec(ctx)
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}