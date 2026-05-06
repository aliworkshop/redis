package redis

import (
	"context"
	"fmt"
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
	"time"
)

func BenchmarkRedisWatch_Concurrent(b *testing.B) {
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

	key := "portfolio:wallet:1:USDT"
	balanceKey := "portfolio:balance:1:USDT"
	blockedKey := "portfolio:blocked:1:USDT"

	// Initialize test keys
	require.NoError(b, redis.Store(ctx, balanceKey, "1000", 0))
	require.NoError(b, redis.Store(ctx, blockedKey, "900", 0))

	b.ResetTimer()

	var wg sync.WaitGroup
	concurrency := 90
	amount := decimal.NewFromInt(10)
	//redis.Unlock(ctx, key)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer func() {
				redis.Unlock(ctx, key)
				wg.Done()
			}()
			for {
				ok, _ := redis.Lock(ctx, key, 10*time.Second)
				if !ok {
					//time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}

			err = redis.Watch(ctx, key, func(db dbcore.Cache) errors.ErrorModel {
				var balance, blocked decimal.Decimal
				if err = db.Load(ctx, balanceKey, &balance); err != nil {
					return err
				}
				if err = db.Load(ctx, blockedKey, &blocked); err != nil {
					return err
				}

				if blocked.LessThan(amount) {
					return errors.Validation().WithMessage("insufficient blocked amount")
				}
				return db.Store(ctx, blockedKey, blocked.Sub(amount).String())
			})

			if err != nil {
				if err.Detail() == "redis: transaction failed" {
					b.Logf("goroutine %d: transaction failed", id)
				} else {
					b.Logf("goroutine %d: other error: %v", id, err)
				}
			}
		}(i)
	}

	wg.Wait()
	var balance, blocked decimal.Decimal
	require.NoError(b, redis.Load(ctx, balanceKey, &balance))
	fmt.Println("balance:", balance)
	require.NoError(b, redis.Load(ctx, blockedKey, &blocked))
	fmt.Println("blocked:", blocked)
	require.Equal(b, decimal.NewFromInt(1000), balance)
	require.Equal(b, decimal.NewFromInt(0), blocked)

	b.StopTimer()
	time.Sleep(1 * time.Second) // allow logs to flush
}