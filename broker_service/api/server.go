package api

import (
	"github.com/RoyceAzure/sexy_gpt/broker_service/repository/account_service_dao"
)

type Server struct {
	AccountServiceDao accountservicedao.IAccountServiceDao
}
