package service

import (
	"context"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ISessionService interface {
	/*
		des: 檢查User session狀態是否正常
		parm: usersId
		errors:
			ErrInvalidSession : session不存在, session過期, session被block, 需要重新登入
			ErrInternal : 未預期錯誤
	*/
	IsUserLogin(ctx context.Context, userId uuid.UUID) (bool, error)
}

func (userService *Service) IsUserLogin(ctx context.Context, userId uuid.UUID) (bool, error) {
	session, err := userService.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: userId,
		Valid: true,
	})
	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			return false, gpt_error.ErrInvalidSession
		}
		return false, gpt_error.ErrInternal
	}

	if time.Now().UTC().After(session.ExpiredAt.Time) || session.IsBlocked {
		_, err := userService.dao.DeleteSession(ctx, session.ID)
		if err != nil {
			return false, gpt_error.ErrInternal
		}
		return false, gpt_error.ErrInvalidSession
	}

	return true, nil
}
