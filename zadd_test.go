package redis

import (
	"context"
	"fmt"
	"github.com/aliworkshop/configer"
	"os"
	"testing"
	"time"
)

type OHLC struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

func TestRedisZAdd(t *testing.T) {
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
	for range 5 {
		time.Sleep(time.Second)
		ohlc := OHLC{
			Timestamp: time.Now().Unix(),
			Open:      50000.0,
			High:      50500.0,
			Low:       49800.0,
			Close:     50400.0,
			Volume:    123.45,
		}
		err = redis.ZAdd(ctx, "ohlc:BTCUSDT:1m", float64(ohlc.Timestamp), ohlc)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("before:")
	res, err := redis.ZRangeLoad(ctx, "ohlc:BTCUSDT:1m", "0", "9999999999", new(OHLC), 0, 100)
	if err != nil {
		panic(err)
	}
	for _, item := range res.([]*OHLC) {
		fmt.Println(item)
	}

	err = redis.ZRemoveByRank(ctx, "ohlc:BTCUSDT:1m", 0, -5) // for keeping latest 4 records
	if err != nil {
		panic(err)
	}

	fmt.Println("after:")
	res, err = redis.ZRangeLoad(ctx, "ohlc:BTCUSDT:1m", "0", "9999999999", new(OHLC), 0, 100)
	if err != nil {
		panic(err)
	}
	for _, item := range res.([]*OHLC) {
		fmt.Println(item)
	}
}