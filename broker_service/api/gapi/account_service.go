package gapi

import (
	"github.com/RoyceAzure/sexy_gpt/broker_service/repository/account_service_dao"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	authorizer        IAuthorizer
	accountServiceDao accountservicedao.IAccountServiceDao
}

func NewAccountServer(authorizer IAuthorizer, accountServiceDao accountservicedao.IAccountServiceDao) *AccountServer {
	return &AccountServer{
		authorizer:        authorizer,
		accountServiceDao: accountServiceDao,
	}
}
