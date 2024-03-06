package pawaiservicedaogo

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
)

func (dao *PawAIServiceDao) FreeChat(ctx context.Context, req *pb.FreeChatRequest, dto model.UserDTO) (*pb.ChatResponse, error) {
	newCtx := util.NewOutGoingMetaDataUserDTO(ctx, dto)
	res, err := dao.client.FreeChat(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (dao *PawAIServiceDao) Chat(ctx context.Context, req *pb.ChatRequest, dto model.UserDTO) (*pb.ChatResponse, error) {
	newCtx := util.NewOutGoingMetaDataUserDTO(ctx, dto)
	res, err := dao.client.Chat(newCtx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
