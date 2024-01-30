package db

import (
	"context"
	"testing"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateSession(t *testing.T) Session {
	user := CreateRandomUsesr(t)
	refreshToken := random.RandomString(10)
	userAgent := random.RandomString(5)
	clientIp := random.RandomString(5)
	session, err := testDao.CreateSession(context.Background(), CreateSessionParams{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
	})
	require.NotEmpty(t, session)
	require.NoError(t, err)
	return session
}
func TestCreateSession(t *testing.T) {
	CreateSession(t)
}

func TestDeleteSession(t *testing.T) {
	session := CreateSession(t)
	deletedSession, err := testDao.DeleteSession(context.Background(), session.ID)
	require.NotEmpty(t, deletedSession)
	require.NoError(t, err)
}

func TestGetSessionByUserId(t *testing.T) {
	session := CreateSession(t)
	getSeession, err := testDao.GetSessionByUserId(context.Background(), session.UserID)
	require.NotEmpty(t, getSeession)
	require.NoError(t, err)
	require.Equal(t, session.UserID, getSeession.UserID)
	require.Equal(t, session.RefreshToken, getSeession.RefreshToken)
	require.Equal(t, session.UserAgent, getSeession.UserAgent)
	require.Equal(t, session.ClientIp, getSeession.ClientIp)
}

func TestGetSessions(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateSession(t)
	}
	getSessions, err := testDao.GetSessions(context.Background(), GetSessionsParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, getSessions, 5)
}

func TestUpdateSession(t *testing.T) {
	session := CreateSession(t)
	arg := UpdateSessionParams{
		IsBlocked: pgtype.Bool{
			Valid: true,
			Bool:  true,
		},
		ID: session.ID,
	}
	_, err := testDao.UpdateSession(context.Background(), arg)
	require.NoError(t, err)
	updatedSession, err := testDao.GetSessionByUserId(context.Background(), session.UserID)

	require.NoError(t, err)
	require.Equal(t, arg.IsBlocked.Bool, updatedSession.IsBlocked)
	require.Equal(t, arg.ID, updatedSession.ID)
	// require.Equal(t, arg.UpDate.Time, updatedRole.UpDate)
}
