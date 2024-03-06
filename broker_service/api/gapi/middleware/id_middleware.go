package middleware

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func IdMiddleWare(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	reqID := uuid.New().String()
	ctx = context.WithValue(ctx, util.RequestIDKey, reqID)
	return handler(ctx, req)
}
