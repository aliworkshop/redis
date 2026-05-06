package redis

import (
	"context"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
)

func (r *repo) Publish(ctx context.Context, channel string, message any) errors.ErrorModel {
	return errors.HandleError(r.client.Publish(ctx, channel, message).Err())
}

func (r *repo) Subscribe(ctx context.Context, channels ...string) <-chan *redis.Message {
	return r.client.Subscribe(ctx, channels...).Channel()
}