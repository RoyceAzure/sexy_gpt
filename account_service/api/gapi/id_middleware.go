package gapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	// logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/remote_dao/logger_distributor"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func IdMiddleWare(ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	reqID := uuid.New().String()
	ctx = context.WithValue(ctx, util.RequestIDKey, reqID)
	return handler(ctx, req)
}

func IdMiddleWareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// log := logger.Logger.Info()
		// log.Msg("in IdMiddleWareHandler")

		if val := req.Header.Get(string(util.RequestIDKey)); val == "" {
			reqID := uuid.New().String()
			req.Header.Set(string(util.RequestIDKey), reqID)
		}

		handler.ServeHTTP(res, req)
	})
}

func CustomMatcher(ctx context.Context, req *http.Request) metadata.MD {
	// 創建一個空的 metadata.MD 對象
	md := metadata.MD{}

	// 從 HTTP 請求中提取 header 並添加到 metadata.MD
	if val := req.Header.Get(string(util.RequestIDKey)); val != "" {
		md[strings.ToLower(string(util.RequestIDKey))] = []string{val}
	}
	return md
}
