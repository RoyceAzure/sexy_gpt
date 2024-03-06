package gapi

import (
	"context"
	"fmt"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/model"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func processChatResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.ChatResponse, error) {
	res := pb.ChatResponse{Ans: msg}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

// 不適用  後續需要修改   因為錯誤訊息沒有統一
func processInitSessionResponse(ctx context.Context, code codes.Code, msg string, codeInternal int32, err error) (*pb.InitSessionResponse, error) {
	res := pb.InitSessionResponse{Code: codeInternal}
	if code != codes.OK {
		if err == nil {
			err = fmt.Errorf(msg)
		}
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

func (server *PawAIServerServer) FreeChat(ctx context.Context, req *pb.FreeChatRequest) (*pb.ChatResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processChatResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	authRes, err := server.accountServiceDao.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: token,
	})

	if err != nil {
		return processChatResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	usetDTO := model.UserDTO{
		UserId:   authRes.User.UserId,
		UserName: authRes.User.UserName,
		Email:    authRes.User.Email,
		RoleName: authRes.User.RoleName,
		RoleId:   authRes.User.RoleId,
	}

	return server.pawaiServiceDao.FreeChat(ctx, req, usetDTO)
}
func (server *PawAIServerServer) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processChatResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	authRes, err := server.accountServiceDao.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: token,
	})

	if err != nil {
		return processChatResponse(ctx, codes.Unauthenticated, err.Error(), err)
	}

	usetDTO := model.UserDTO{
		UserId:   authRes.User.UserId,
		UserName: authRes.User.UserName,
		Email:    authRes.User.Email,
		RoleName: authRes.User.RoleName,
		RoleId:   authRes.User.RoleId,
	}

	return server.pawaiServiceDao.Chat(ctx, req, usetDTO)
}
func (server *PawAIServerServer) InitSession(ctx context.Context, req *pb.InitSessionRequest) (*pb.InitSessionResponse, error) {
	_, token, err := server.authorizer.AuthorizToken(ctx)
	if err != nil {
		return processInitSessionResponse(ctx, codes.Unauthenticated, err.Error(), 500, err)
	}

	authRes, err := server.accountServiceDao.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: token,
	})

	if err != nil {
		return processInitSessionResponse(ctx, codes.Unauthenticated, err.Error(), 401, err)
	}

	usetDTO := model.UserDTO{
		UserId:   authRes.User.UserId,
		UserName: authRes.User.UserName,
		Email:    authRes.User.Email,
		RoleName: authRes.User.RoleName,
		RoleId:   authRes.User.RoleId,
	}

	return server.pawaiServiceDao.InitSession(ctx, req, usetDTO)
}
