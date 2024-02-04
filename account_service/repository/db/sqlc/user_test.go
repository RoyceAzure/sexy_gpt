package db

import (
	"context"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUsesr(t)
}

func CreateRandomUsesr(t *testing.T) User {
	userName := random.RandomString(10)
	email := userName + "@gmail.com"
	password := random.RandomString(15)
	hashPaw, err := util.HashPassword(password)
	require.NoError(t, err)

	user, err := testDao.CreateUser(context.Background(), CreateUserParams{
		UserName:       userName,
		Email:          email,
		HashedPassword: hashPaw,
		CrUser:         "SYSTEM",
	})
	require.NotEmpty(t, user)
	require.NoError(t, err)
	return user
}

func TestGetUserByEmail(t *testing.T) {
	user := CreateRandomUsesr(t)
	returnUser, err := testDao.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.EqualValues(t, user, returnUser)
}

func TestGetUsers(t *testing.T) {
	CreateRandomUsesr(t)
	CreateRandomUsesr(t)
	users, err := testDao.GetUsers(context.Background(), GetUsersParams{
		Limit:  2,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestUpdateUser(t *testing.T) {
	user := CreateRandomUsesr(t)
	newName := random.RandomString(5)
	newEmail := newName + "@gmail.com"
	newPas := random.RandomString(10)
	newHashPas, err := util.HashPassword(newPas)
	require.NoError(t, err)

	now := time.Now().UTC()
	arg := UpdateUserParams{
		UserName: pgtype.Text{
			Valid:  true,
			String: newName,
		},
		Email: pgtype.Text{
			Valid:  true,
			String: newEmail,
		},
		IsEmailVerified: pgtype.Bool{
			Valid: true,
			Bool:  true,
		},
		HashedPassword: pgtype.Text{
			Valid:  true,
			String: newHashPas,
		},
		PasswordChangedAt: pgtype.Timestamptz{
			Valid: true,
			Time:  now,
		},
		IsInternal: pgtype.Bool{
			Valid: true,
			Bool:  true,
		},
		UpDate: pgtype.Timestamptz{
			Valid: true,
			Time:  now,
		},
		UpUser: pgtype.Text{
			Valid:  true,
			String: "testing",
		},
		UserID: user.UserID,
	}
	_, err = testDao.UpdateUser(context.Background(), arg)

	require.NoError(t, err)

	updatedUser, err := testDao.GetUser(context.Background(), user.UserID)
	require.NoError(t, err)

	require.Equal(t, arg.UserName.String, updatedUser.UserName)
	require.Equal(t, arg.Email.String, updatedUser.Email)
	require.Equal(t, arg.IsEmailVerified.Bool, updatedUser.IsEmailVerified)
	require.Equal(t, arg.HashedPassword.String, updatedUser.HashedPassword)
	require.Equal(t, arg.UpUser.String, updatedUser.UpUser.String)
	// require.Equal(t, arg.PasswordChangedAt.Time, updatedUser.PasswordChangedAt.UTC())
	// require.Equal(t, arg.UpDate.Time, updatedUser.UpDate.Time)
	require.Equal(t, arg.IsInternal.Bool, updatedUser.IsInternal)
}
