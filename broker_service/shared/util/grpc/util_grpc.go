package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CustomrErrHttpHeader(ctx context.Context, header metadata.MD) {
	// var headerMD metadata.MD
	// headerValues := header.Get(util.DBMSGKey)
	// if len(headerValues) > 0 {
	// 	headerMD = metadata.Pairs(
	// 		util.DBMSGKey, headerValues[0],
	// 	)
	// } else {
	// 	headerMD = metadata.Pairs(
	// 		util.DBMSGKey, err.Error(),
	// 	)
	// }
	grpc.SendHeader(ctx, header)
}
