package worker

import (
	"context"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/mail"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	MailQueue            = "mailQueue"
	SyncStockQueue       = "syncStockQueue"
	StockTransationQueue = "stockTransationQueue"
)

/*
asynq有定義 processor 的interface
*/
type TaskProcessor interface {
	Start() error
	StartWithHandler(handler *asynq.ServeMux) error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	dao    db.Dao
	mailer mail.EmailSender
}

/*
server是processot自己建立的? 是host建立asynq server去遠端redis 拿取任務資料
*/
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt,
	dao db.Dao,
	mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				StockTransationQueue: 10,
				MailQueue:            8,
				SyncStockQueue:       6,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("body", task.Payload()).
					Msg("process task failed")
			}),
			Logger: NewLoggerAdapter(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		dao:    dao,
		mailer: mailer,
	}
}

/*
案照這個設計思維   所有processor應可以連到同一個redis server
且每個processor start 就是指start一個路由與handler的對應  且會永久執行?  對  就跟http server很像
一旦執行後就會被block  所以要用go routine執行
asynq 本身processor處理每個task就是使用go routine處理
error是指 sersver startup 過程的錯誤
*/
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
func (processor *RedisTaskProcessor) StartWithHandler(handler *asynq.ServeMux) error {
	return processor.server.Start(handler)
}
