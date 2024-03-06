package accountservicedao

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func (dao *AccountServiceDao) Login(ctx context.Context, req *pb.LoginRequset) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.Login(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) Logout(ctx context.Context, req *pb.LogoutRequset, accessToken string) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.Logout(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequset, accessToken string) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	res, err := dao.client.RefreshToken(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) SSOGoogleLogin(ctx context.Context, req *pb.GoogleIDTokenRequest) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.SSOGoogleLogin(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *AccountServiceDao) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.ValidateToken(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
