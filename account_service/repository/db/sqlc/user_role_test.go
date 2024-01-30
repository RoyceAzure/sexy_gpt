package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateUserRole(t *testing.T) UserRole {
	user := CreateRandomUsesr(t)
	role := CreateRole(t, "")

	usesrRole, err := testDao.CreateUserRole(context.Background(), CreateUserRoleParams{
		UserID: user.UserID,
		RoleID: role.RoleID,
		CrUser: "SYSTEM",
	})
	require.NotEmpty(t, usesrRole)
	require.NoError(t, err)
	return usesrRole
}
func TestCreateUserRole(t *testing.T) {
	CreateRole(t, "")
}

func TestGetUserRoleByUserId(t *testing.T) {
	userRole := CreateUserRole(t)
	getUserRoles, err := testDao.GetUserRoleByUserId(context.Background(), userRole.UserID)
	require.NotEmpty(t, getUserRoles)
	require.NoError(t, err)
	for _, userRole := range getUserRoles {
		require.Equal(t, userRole.UserID, userRole.UserID)
	}
}

func TestGetUserRoles(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateUserRole(t)
	}
	userRoles, err := testDao.GetUserRoles(context.Background(), GetUserRolesParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, userRoles, 5)
}

func TestUpdateUserRole(t *testing.T) {
	userRole := CreateUserRole(t)
	newRole := CreateRole(t, "")

	now := time.Now().UTC()

	arg := UpdateUserRoleParams{
		RoleID: pgtype.UUID{
			Valid: true,
			Bytes: newRole.RoleID.Bytes,
		},
		UpDate: pgtype.Timestamptz{
			Valid: true,
			Time:  now,
		},
		UpUser: pgtype.Text{
			Valid:  true,
			String: "testing",
		},
		UserID: userRole.UserID,
	}
	updatedRole, err := testDao.UpdateUserRole(context.Background(), arg)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, arg.RoleID.Bytes, updatedRole.RoleID.Bytes)
	// require.Equal(t, arg.UpDate.Time, updatedRole.UpDate)
	require.Equal(t, arg.UpUser.String, updatedRole.UpUser.String)
}
