package service

import db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"

type IService interface {
	IUserService
	ISessionService
}

type Service struct {
	dao db.Dao
}

func NewService(dao db.Dao) *Service {
	return &Service{
		dao: dao,
	}
}
