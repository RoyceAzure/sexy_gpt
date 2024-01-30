package gapi

import (
	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/account_service/service"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	config     config.Config
	dao        db.Dao
	tokenMaker token.Maker
	Service    service.IService
}

func NewServer(config config.Config, dao db.Dao, tokenMaker token.Maker, service service.IService) (*Server, error) {
	server := &Server{
		config:     config,
		dao:        dao,
		tokenMaker: tokenMaker,
		Service:    service,
	}
	return server, nil
}

func (server *Server) HandleAPIError(code codes.Code, err error, msg string) error {
	if logger.Logger != nil {
		logger.Logger.Error().Err(err).Msg(msg)
	}
	return status.Errorf(code, msg)
}
