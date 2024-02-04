package worker

import (
	"context"
	"encoding/json"
	"fmt"

	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserId     uuid.UUID `json:"user_id"`
	SecretCode string    `json:"secret_code"`
}

/*
製作task  這裡把所有資訊都包進task裡面  包括retry delay 甚至是要使用甚麼優先級的queue

使用client enqueue task
*/
func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	md := util.ExtractMetaData(ctx)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Any("meta", md).
			Msg("enqueued task failed")
		return fmt.Errorf("failed to marshal task payload %w", err)
	}

	//沒有提到要如何執行task , 只有給payload跟設定
	task := asynq.NewTask(TaskSendVerifyEmail, jsonData, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("type", task.Type()).
			Any("meta", md).
			Bytes("body", task.Payload()).
			Str("queue", taskInfo.Queue).
			Int("max_retry", taskInfo.MaxRetry).
			Msg("enqueued task")
		return fmt.Errorf("failed to enqueue task %w", err)
	}

	logger.Logger.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueued task")
	return nil
}

/*
因為這個func 簽名(ctx context.Context, task *asynq.Task) error 是asynq定義
所以asynq內部會根據error內容作相應處理

如果回傳得error 有wrap asynq.SkipRetry，則不會執行retry

注意這段handler是在redis server內部執行  所以它回傳的err要額外設置  我們才能看得見

err == sql.ErrNORows 有可能是db 還沒有完成commit，所以就讓他retry
*/
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	md := util.ExtractMetaData(ctx)
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload %w", asynq.SkipRetry)
	}
	user, err := processor.dao.GetUser(ctx, pgtype.UUID{
		Bytes: payload.UserId,
		Valid: true,
	})
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user doesn't exists %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get usesr %w", asynq.SkipRetry)
	}

	//TODO : send email to user
	subject := "Welcome to Sexy GPT"
	verifyURL := fmt.Sprintf("http://localhost:8081/v1/vertifyEmail?user_id=%s&secret_code=%s",
		payload.UserId, payload.SecretCode)

	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering <br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.UserName, verifyURL)

	to := []string{user.Email}
	bcc := []string{"roycewnag@gmail.com"}

	err = processor.mailer.SendEmail(subject, content, to, nil, bcc, nil)
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}

	logger.Logger.Info().Str("type", task.Type()).
		Bytes("task payload", task.Payload()).
		Any("meta", md).
		Str("user email", user.Email).
		Msg("porcessed task")
	return nil
}
