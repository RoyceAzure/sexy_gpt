package pawaiservicedaogo

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func (dao *PawAIServiceDao) InitSession(ctx context.Context, req *pb.InitSessionRequest, dto model.UserDTO) (*pb.InitSessionResponse, error) {
	newCtx := util.NewOutGoingMetaDataUserDTO(ctx, dto)
	res, err := dao.client.InitSession(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
