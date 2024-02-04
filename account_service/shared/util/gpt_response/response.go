package gpt_response

import (
	"context"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type IProteMessage interface {
// 	Reset()
// 	String() string
// 	ProtoMessage()
// }

/*
ProcessResponse 處理res, 記錄log，并返回pb.Response和错误。
code grpc code。
msg 響應訊息。
err 日誌額外季紀錄的錯誤訊息
return *pbany.Any
*/
func ProcessResponse(ctx context.Context, code codes.Code, msg string, err error) (*pb.Response, error) {
	res := pb.Response{Message: msg}
	if code != codes.OK {
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}

/*
ProcessResponse 處理res, 記錄log，并返回pb.Response和错误。
code grpc code。
msg 響應訊息。
err 日誌額外季紀錄的錯誤訊息。
return []*pbany.Any
*/
func ProcessResponses(ctx context.Context, code codes.Code, msg string, err error) (*pb.Responses, error) {
	res := pb.Responses{Message: msg}
	if code != codes.OK {
		util.NewOutGoingMetaDataKV(ctx, util.DBMSGKey, err.Error())
		return &res, status.Errorf(code, msg)
	}
	return &res, nil
}
