package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"google.golang.org/grpc/metadata"
)

func CustomMatcher(ctx context.Context, req *http.Request) metadata.MD {
	// 創建一個空的 metadata.MD 對象
	md := metadata.MD{}

	// 從 HTTP 請求中提取 header 並添加到 metadata.MD
	if val := req.Header.Get(string(util.RequestIDKey)); val != "" {
		md[strings.ToLower(string(util.RequestIDKey))] = []string{val}
	}
	return md
}
