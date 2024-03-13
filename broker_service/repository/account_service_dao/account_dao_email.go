package accountservicedao

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	util_grpc "github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (dao *AccountServiceDao) SendVertifyEmai(ctx context.Context, req *pb.SendVertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	var header metadata.MD
	res, err := dao.client.SendVertifyEmai(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) VertifyEmail(ctx context.Context, req *pb.VertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	var header metadata.MD
	res, err := dao.client.VertifyEmail(newCtx, req, grpc.Header(&header))
	if err != nil {
		util_grpc.CustomrErrHttpHeader(ctx, header)
		return nil, err
	}
	return res, nil
}
