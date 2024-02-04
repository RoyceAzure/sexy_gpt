package gapi

import (
	"context"
	"net/http"
	"time"

	logger "github.com/RoyceAzure/sexy_gpt/account_service/repository/logger_distributor"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util"
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

/*
使用裝飾者模式來擴展ResponseWriter的功能
多了紀錄StatusCode
*/
type ResponseRecoder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecoder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecoder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

/*
當您呼叫 http.HandlerFunc(MyHandler) 時，您其實是在進行一個型別轉換，
將 MyHandler 這個函數轉換為 HandlerFunc 型別
這裡使用匿名函數自訂func

	func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
		f(w, r)
	}

當呼叫ServeHTTP時  實際上就是執行此type的func
*/
func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log := logger.Logger.Info()
		// log.Msg("in HttpLogger")
		startTime := time.Now()

		rec := &ResponseRecoder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		reqId := req.Header.Get(string(util.RequestIDKey))

		handler.ServeHTTP(rec, req)

		if rec.StatusCode != http.StatusOK {
			log = logger.Logger.Error().Bytes("res body", rec.Body)
		}

		dbMsg := rec.Header().Get(util.DBMSGKey)

		duration := time.Since(startTime).Milliseconds()
		log.Str("protocol", "http").
			Str("method", req.Method).
			Str("request_id", reqId).
			Str(util.DBMSGKey, dbMsg).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Int64("duration in ms", duration).
			Msg("received a HTTP request")
	})
}
