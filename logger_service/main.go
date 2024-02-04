package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RoyceAzure/sexy_gpt/logger_service/asynq_worker"
	repository "github.com/RoyceAzure/sexy_gpt/logger_service/repository/mongodb"
	logservice "github.com/RoyceAzure/sexy_gpt/logger_service/service"
	"github.com/RoyceAzure/sexy_gpt/logger_service/shared/util/config"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	config, err := config.LoadConfig(".") //表示讀取當前資料夾
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := repository.ConnectToMongo(ctx, config.MongodbAddress)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot connect to mongo db")
	}

	err = repository.InitCollection(client)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("init mongo db failed")
	}

	mongoDao := repository.NewMongoDao(client)

	mongoLogger := logservice.NewMongoLogger(mongoDao)
	err = logservice.SetUpMutiMongoLogger(mongoLogger, config.ServiceID)
	if err != nil {
		logservice.Logger.Fatal().
			Err(err).
			Msg("cannot connect to mongo db")
	}
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisQueueAddress,
	}
	processer := asynq_worker.NewRedisTaskProcessor(redisOpt, mongoDao)

	logservice.Logger.Info().Msg("Start logger service")

	go processer.Start()

	sigChan := make(chan os.Signal, 1)
	// 监听所有中断信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞，直到从 sigChan 接收到信号
	sig := <-sigChan
	logservice.Logger.Info().Msgf("Received signal: %s", sig)

	// ... 在这里可以进行优雅的关闭处理，比如关闭数据库连接，清理资源等

	// 停止服务（如果 Start 是非阻塞的）

	logservice.Logger.Info().Msg("Stop logger service")
	processer.Stop()
}
