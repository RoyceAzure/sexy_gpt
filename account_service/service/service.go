package service

import (
	db "github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc"
	"github.com/RoyceAzure/sexy_gpt/account_service/worker"
)

type IService interface {
	IUserService
	ISessionService
	IVertifyEmailService
}

type Service struct {
	dao         db.Dao
	asynqWorker worker.ITaskDistributor
}

func NewService(dao db.Dao, asynqWorker worker.ITaskDistributor) IService {
	return &Service{
		dao:         dao,
		asynqWorker: asynqWorker,
	}
}
