package middleware

import (
	"context"
	"time"

	logger "github.com/RoyceAzure/sexy_gpt/broker_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	startTime := time.Now()
	log := logger.Logger.Info()

	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	if err != nil {
		log = logger.Logger.Error().Err(err)
	}

	mtda := util.ExtractMetaData(ctx)

	log.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Any("meta", mtda).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")
	return result, err
}
