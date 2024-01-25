package gapi

import (
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	config config.Config
	dao    db.Dao
}

func NewServer(config config.Config, dao db.Dao) (*Server, error) {
	server := &Server{
		config: config,
		dao:    dao,
	}
	return server, nil
}
