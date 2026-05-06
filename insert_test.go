package redis

import (
	"context"
	"github.com/aliworkshop/configer"
	"os"
	"testing"
)

func TestRedisInsert(t *testing.T) {
	registry := configer.New()
	registry.SetConfigType("yaml")
	f, e := os.Open("./config.sample.yaml")
	if e != nil {
		panic("cannot read config: " + e.Error())
	}
	e = registry.ReadConfig(f)
	if e != nil {
		panic("cannot read config" + e.Error())
	}

	redis := NewRedisRepository(registry)
	err := redis.Initialize()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = redis.Store(ctx, "test", "value")
	if err != nil {
		panic(err)
	}
}