package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ISessionService interface {
	/*
		des:
			檢查User session狀態是否正常
		parm:
			usersId
		errors:
			ErrInvalidSession : session不存在, session過期, session被block, 需要重新登入
			ErrInternal : 未預期錯誤
	*/
	IsUserLogin(ctx context.Context, userId uuid.UUID) (bool, error)
	LoginCreateSession(ctx context.Context, userId uuid.UUID, tokenMaker token.Maker, config config.Config) (string, error)
}

func (sessionService *Service) IsUserLogin(ctx context.Context, userId uuid.UUID) (bool, error) {
	session, err := sessionService.dao.GetSessionByUserId(ctx, pgtype.UUID{
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
		_, err := sessionService.dao.DeleteSession(ctx, session.ID)
		if err != nil {
			return false, gpt_error.ErrInternal
		}
		return false, gpt_error.ErrInvalidSession
	}

	return true, nil
}

// create session via user login
// if session not exists, create session
// if session expired, delete session than create
// if session create succedded, or session has not expired, return refreshtoken of session
// erro : only internal err
func (sessionService *Service) LoginCreateSession(ctx context.Context, userId uuid.UUID, tokenMaker token.Maker, config config.Config) (string, error) {

	oldSession, err := sessionService.dao.GetSessionByUserId(ctx, pgtype.UUID{
		Bytes: userId,
		Valid: true,
	})

	var refreshToken string
	var refreshPayLoad *token.TokenPayload

	if err != nil {
		if err.Error() == gpt_error.DB_ERR_NOT_FOUND.Error() {
			refreshToken, refreshPayLoad, err = tokenMaker.CreateToken(nil, "refresh", config.AUTH_ISSUER, config.RefreshTokenDuration)
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}

			_, err = sessionService.dao.CreateSession(ctx, db.CreateSessionParams{
				UserID: pgtype.UUID{
					Bytes: userId,
					Valid: true,
				},
				RefreshToken: refreshToken,
				UserAgent:    "todo",
				ClientIp:     "todo",
				ExpiredAt: pgtype.Timestamptz{
					Time:  refreshPayLoad.ExpiredAt,
					Valid: true,
				},
			})
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}
		} else {
			_, err := sessionService.dao.DeleteSession(ctx, oldSession.ID)
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}
		}
	} else {
		if time.Now().After(oldSession.ExpiredAt.Time) {
			_, err := sessionService.dao.DeleteSession(ctx, oldSession.ID)
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}
			refreshToken, refreshPayLoad, err := tokenMaker.CreateToken(nil, "refresh", config.AUTH_ISSUER, config.RefreshTokenDuration)
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}

			_, err = sessionService.dao.CreateSession(ctx, db.CreateSessionParams{
				UserID: pgtype.UUID{
					Bytes: userId,
					Valid: true,
				},
				RefreshToken: refreshToken,
				UserAgent:    "todo",
				ClientIp:     "todo",
				ExpiredAt: pgtype.Timestamptz{
					Time:  refreshPayLoad.ExpiredAt,
					Valid: true,
				},
			})
			if err != nil {
				return "", fmt.Errorf("err : %s, %w", err.Error(), gpt_error.ErrInternal)
			}
		}
		refreshToken = oldSession.RefreshToken
	}
	return refreshToken, nil
}
