package middleware

type MiddlewareResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

// for call accountService, 只需檢查token 格式, 驗證使用者交給account service
// func HttpAuthMiddleWare(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var res MiddlewareResponse
// 		tokenMaker, err := token.GetSingleTonPasetoTokenMaker()
// 		if err != nil {
// 			w.WriteHeader(500)
// 			res.Message = "Internal err"
// 			return
// 		}

// 	})
// }
