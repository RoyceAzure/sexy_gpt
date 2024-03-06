package pawaiservicedaogo

import (
	"context"
	"fmt"
	"time"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/gpt_error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IPawAIServiceDao interface {
	FreeChat(context.Context, *pb.FreeChatRequest, model.UserDTO) (*pb.ChatResponse, error)
	Chat(context.Context, *pb.ChatRequest, model.UserDTO) (*pb.ChatResponse, error)
	InitSession(context.Context, *pb.InitSessionRequest, model.UserDTO) (*pb.InitSessionResponse, error)
}

type PawAIServiceDao struct {
	client pb.PawAIServiceClient
}

func NewPawAIServiceDao(address string) (IPawAIServiceDao, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("can't connect grpc server, %w", gpt_error.ErrInternal)
	}

	client := pb.NewPawAIServiceClient(conn)

	return &PawAIServiceDao{
		client: client,
	}, func() { conn.Close() }, nil
}
