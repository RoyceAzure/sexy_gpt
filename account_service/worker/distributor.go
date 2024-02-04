package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type ITaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

/*
注意  asynq.Client 並沒有提到是for redis, 看來是通用的??

A Client is responsible for scheduling tasks.

A Client is used to register tasks that should be processed immediately or some time in the future.

Clients are safe for concurrent use by multiple goroutines
*/
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) ITaskDistributor {
	//asynq.NewClient() 這裡就有提for redis
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
