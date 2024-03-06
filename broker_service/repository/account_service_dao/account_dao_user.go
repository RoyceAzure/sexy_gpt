package accountservicedao

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func (dao *AccountServiceDao) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.CreateUser(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) GetUser(ctx context.Context, req *pb.GetUserRequest, accessToken string) (*pb.UserDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.GetUser(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) GetUsers(ctx context.Context, req *pb.GetUsersRequest, accessToken string) (*pb.UserDTOsResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.GetUsers(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest, accessToken string) (*pb.UserDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.GetUserByEmail(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, accessToken string) (*pb.UserDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.UpdateUser(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) UpdateUserPas(ctx context.Context, req *pb.UpdateUserPasRequest, accessToken string) (*pb.UserDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.UpdateUserPas(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
