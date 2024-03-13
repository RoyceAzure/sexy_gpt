package accountservicedao

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	util_grpc "github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (dao *AccountServiceDao) Login(ctx context.Context, req *pb.LoginRequset) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	var header metadata.MD
	res, err := dao.client.Login(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) Logout(ctx context.Context, req *pb.LogoutRequset, accessToken string) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	var header metadata.MD
	res, err := dao.client.Logout(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequset, accessToken string) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, accessToken)
	var header metadata.MD
	res, err := dao.client.RefreshToken(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) SSOGoogleLogin(ctx context.Context, req *pb.GoogleIDTokenRequest) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	var header metadata.MD
	res, err := dao.client.SSOGoogleLogin(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}

func (dao *AccountServiceDao) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.AuthDTOResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	var header metadata.MD
	res, err := dao.client.ValidateToken(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
