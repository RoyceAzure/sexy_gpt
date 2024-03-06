package gapi

import (
	accountservicedao "github.com/RoyceAzure/sexy_gpt/broker_service/repository/account_service_dao"
	"github.com/RoyceAzure/sexy_gpt/broker_service/repository/pawai_service_dao"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
)

type PawAIServerServer struct {
	pb.UnimplementedPawAIServiceServer
	authorizer        IAuthorizer
	pawaiServiceDao   pawaiservicedaogo.IPawAIServiceDao
	accountServiceDao accountservicedao.IAccountServiceDao
}

func NewPawAIServerServer(
	authorizer IAuthorizer,
	pawaiServiceDao pawaiservicedaogo.IPawAIServiceDao,
	accountServiceDao accountservicedao.IAccountServiceDao,
) *PawAIServerServer {
	return &PawAIServerServer{
		authorizer:        authorizer,
		pawaiServiceDao:   pawaiServiceDao,
		accountServiceDao: accountServiceDao,
	}
}
