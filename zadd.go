package redis

import (
	"context"
	"encoding/json"
	"github.com/aliworkshop/errors"
	"github.com/redis/go-redis/v9"
	"reflect"
)

func (r *repo) ZAdd(ctx context.Context, key string, score float64, member any) errors.ErrorModel {
	data, err := json.Marshal(member)
	if err != nil {
		return errors.Internal(err)
	}
	err = r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: string(data),
	}).Err()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) ZRangeLoad(ctx context.Context, key string, min, max string, result any, offset, count int64) (any, errors.ErrorModel) {
	res, err := r.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
	if err != nil {
		return nil, errors.Internal(err)
	}

	typ := reflect.TypeOf(result)
	slice := reflect.New(reflect.SliceOf(typ)).Elem()
	for _, item := range res {
		elm := reflect.New(typ.Elem())
		if err = json.Unmarshal([]byte(item), elm.Interface()); err != nil {
			return nil, errors.Internal(err)
		}
		slice = reflect.Append(slice, elm)
	}
	return slice.Interface(), nil
}

func (r *repo) ZRemove(ctx context.Context, key string, min, max string) errors.ErrorModel {
	_, err := r.client.ZRemRangeByScore(ctx, key, min, max).Result()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) ZRemoveByRank(ctx context.Context, key string, from, to int64) errors.ErrorModel {
	_, err := r.client.ZRemRangeByRank(ctx, key, from, to).Result()
	if err != nil {
		return errors.Internal(err)
	}
	return nil
}

func (r *repo) ZExist(ctx context.Context, key, member string) (bool, errors.ErrorModel) {
	count, err := r.client.ZCount(ctx, key, member, member).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, errors.Internal(err)
	}
	return count > 0, nil
}

func (r *repo) ZMaxScore(ctx context.Context, key string) (float64, errors.ErrorModel) {
	results, err := r.client.ZRevRangeWithScores(ctx, key, 0, 0).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, errors.Internal(err)
	}
	if len(results) == 0 {
		return 0, nil
	}
	top := results[0]
	return top.Score, nil
}