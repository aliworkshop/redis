package redis

import (
	"context"
	"fmt"
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
	"github.com/shopspring/decimal"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRedisExec(t *testing.T) {
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

	err = redis.Store(ctx, "amount", 1000)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	wg := new(sync.WaitGroup)
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
		trylock:
			ok, _ := redis.Lock(ctx, "test", 10*time.Second)
			if ok {
				err = redis.Watch(ctx, "test", func(db dbcore.Cache) errors.ErrorModel {
					err = db.Store(ctx, "ali", "test")
					if err != nil {
						return err
					}

					amount, err := db.GetInt(ctx, "amount")
					if err != nil {
						return err
					}
					fmt.Println("amount: ", amount)
					if amount < 1 {
						return errors.Validation().WithMessage("Insufficient Balance")
					}

					err = db.DecrBy(ctx, "amount", 1)
					if err != nil {
						return err
					}
					t, err := db.GetExpiration(ctx, "amount")
					if err != nil {
						return err
					}
					fmt.Println(t.String())
					return nil
				})
				if err != nil {
					fmt.Println(time.Since(now))
					panic(err)
				}
				redis.Unlock(ctx, "test")
				wg.Done()
			}
			time.Sleep(10 * time.Millisecond)
			goto trylock
		}()
	}
	wg.Wait()

	var data decimal.Decimal
	err = redis.Load(ctx, "amount", &data)
	if err != nil {
		panic(err)
	}
	fmt.Println("amount:", data.String())
}