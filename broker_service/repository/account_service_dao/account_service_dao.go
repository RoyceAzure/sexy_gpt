package accountservicedao

import (
	"context"
	"fmt"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IAccountServiceDao interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserDTOResponse, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest, accessToken string) (*pb.UserDTOResponse, error)
	GetUsers(ctx context.Context, in *pb.GetUsersRequest, accessToken string) (*pb.UserDTOsResponse, error)
	GetUserByEmail(ctx context.Context, in *pb.GetUserByEmailRequest, accessToken string) (*pb.UserDTOResponse, error)
	UpdateUser(ctx context.Context, in *pb.UpdateUserRequest, accessToken string) (*pb.UserDTOResponse, error)
	UpdateUserPas(ctx context.Context, in *pb.UpdateUserPasRequest, accessToken string) (*pb.UserDTOResponse, error)
	Login(ctx context.Context, in *pb.LoginRequset) (*pb.AuthDTOResponse, error)
	Logout(ctx context.Context, in *pb.LogoutRequset, accessToken string) (*pb.AuthDTOResponse, error)
	RefreshToken(ctx context.Context, in *pb.RefreshTokenRequset, accessToken string) (*pb.AuthDTOResponse, error)
	SendVertifyEmai(ctx context.Context, in *pb.SendVertifyEmailRequset) (*pb.VertifyEmailResponse, error)
	VertifyEmail(ctx context.Context, in *pb.VertifyEmailRequset) (*pb.VertifyEmailResponse, error)
	SSOGoogleLogin(ctx context.Context, in *pb.GoogleIDTokenRequest) (*pb.AuthDTOResponse, error)
	ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.AuthDTOResponse, error)
}

type AccountServiceDao struct {
	client pb.AccountServiceClient
}

func NewAccountServiceDao(address string) (IAccountServiceDao, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("can't connect grpc server, %w", gpt_error.ErrInternal)
	}

	client := pb.NewAccountServiceClient(conn)

	return &AccountServiceDao{
		client: client,
	}, func() { conn.Close() }, nil
}
