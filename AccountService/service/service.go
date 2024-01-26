package service

type IService interface {
	IUserService
}

type Service struct {
	UserService
}
