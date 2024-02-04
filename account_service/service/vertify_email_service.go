package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/RoyceAzure/sexy_gpt/account_service/worker"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

type IVertifyEmailService interface {
	/*
		處理email vertify 舊資料，
		建立新email vertify資料，
		交給task Distrubutor，
		覓等操作，
		errors:
			ErrInternal
	*/
	SendVertifyEmail(ctx context.Context, userId uuid.UUID, email string) error
}

func (s *Service) SendVertifyEmail(ctx context.Context, userId uuid.UUID, email string) error {
	emails, err := s.dao.GetVertifyEmailByEmail(ctx, email)
	if err == nil {
		for _, email := range emails {
			_, err = s.dao.UpdateVertifyEmail(ctx, db.UpdateVertifyEmailParams{
				ID: email.ID,
				IsValidated: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
			})
			if err != nil && err.Error() != gpt_error.DB_ERR_NOT_FOUND.Error() {
				return fmt.Errorf("update vertify email failed, %w", asynq.SkipRetry)
			}
		}
	} else if err.Error() != gpt_error.DB_ERR_NOT_FOUND.Error() {
		return fmt.Errorf("get vertify email by email failed, %w", asynq.SkipRetry)
	}

	secretCode := random.RandomString(32)

	_, err = s.dao.CreateVertifyEmail(ctx, db.CreateVertifyEmailParams{
		UserID: pgtype.UUID{
			Bytes: userId,
			Valid: true,
		},
		Email:      email,
		SecretCode: secretCode,
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		UserId:     userId,
		SecretCode: secretCode,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.MailQueue),
	}
	err = s.asynqWorker.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return fmt.Errorf("%s, %w", err.Error(), gpt_error.ErrInternal)
	}
	return nil
}
