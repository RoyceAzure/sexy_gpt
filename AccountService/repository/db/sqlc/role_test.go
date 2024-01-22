package db

import (
	"context"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateRole(t *testing.T) Role {
	roleName := random.RandomString(5)

	role, err := testDao.CreateRole(context.Background(), CreateRoleParams{
		RoleName: roleName,
		CrUser:   "SYSTEM",
	})
	require.NotEmpty(t, role)
	require.NoError(t, err)
	return role
}
func TestCreateRole(t *testing.T) {
	CreateRole(t)
}

func TestGetRole(t *testing.T) {
	role := CreateRole(t)
	getRole, err := testDao.GetRole(context.Background(), role.RoleID)
	require.NotEmpty(t, getRole)
	require.NoError(t, err)
}

func TestGetRoleByRoleName(t *testing.T) {
	role := CreateRole(t)
	getRole, err := testDao.GetRoleByRoleName(context.Background(), role.RoleName)
	require.NotEmpty(t, getRole)
	require.NoError(t, err)
}

func TestGetRoles(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRole(t)
	}
	roles, err := testDao.GetRoles(context.Background(), GetRolesParams{
		Limit:  5,
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, roles, 5)
}

func TestUpdateRole(t *testing.T) {
	role := CreateRole(t)
	roleName := random.RandomString(5)
	now := time.Now().UTC()
	arg := UpdateRoleParams{
		RoleName: pgtype.Text{
			Valid:  true,
			String: roleName,
		},
		IsEnable: pgtype.Bool{
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
		RoleID: role.RoleID,
	}
	_, err := testDao.UpdateRole(context.Background(), arg)
	require.NoError(t, err)
	updatedRole, err := testDao.GetRole(context.Background(), role.RoleID)

	require.NoError(t, err)
	require.Equal(t, arg.RoleName.String, updatedRole.RoleName)
	require.Equal(t, arg.IsEnable.Bool, updatedRole.IsEnable)
	// require.Equal(t, arg.UpDate.Time, updatedRole.UpDate)
	require.Equal(t, arg.UpUser.String, updatedRole.UpUser.String)
}
