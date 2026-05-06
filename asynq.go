package redis

import (
	"github.com/hibiken/asynq"
)

func (r *repo) InitAsynqServer(queues map[string]int) *asynq.Server {
	queues["low"] = 2
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: r.config.Addr, Password: r.config.Password, DB: r.config.DB},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 25,
			// Optionally specify multiple queues with different priority.
			Queues: queues,
			// See the godoc for other configuration options
		},
	)
}

func (r *repo) InitAsynqScheduler() *asynq.Scheduler {
	srv := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: r.config.Addr, Password: r.config.Password, DB: r.config.DB},
		&asynq.SchedulerOpts{LogLevel: asynq.DebugLevel},
	)
	return srv
}

func (r *repo) InitAsynqClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     r.config.Addr,
		Password: r.config.Password,
		DB:       r.config.DB,
	})
}