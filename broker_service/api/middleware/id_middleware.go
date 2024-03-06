package middleware

import (
	"net/http"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util"
	"github.com/google/uuid"
)

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
