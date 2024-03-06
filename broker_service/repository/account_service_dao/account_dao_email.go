package accountservicedao

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func (dao *AccountServiceDao) SendVertifyEmai(ctx context.Context, req *pb.SendVertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.SendVertifyEmai(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *AccountServiceDao) VertifyEmail(ctx context.Context, req *pb.VertifyEmailRequset) (*pb.VertifyEmailResponse, error) {
	newCtx := util.NewOutGoingMetaData(ctx, "")
	res, err := dao.client.VertifyEmail(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
