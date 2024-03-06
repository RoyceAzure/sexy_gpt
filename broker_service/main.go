package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/RoyceAzure/sexy_gpt/broker_service/api/gapi"
	"github.com/RoyceAzure/sexy_gpt/broker_service/api/middleware"
	"github.com/RoyceAzure/sexy_gpt/broker_service/api/token"
	accountservicedao "github.com/RoyceAzure/sexy_gpt/broker_service/repository/account_service_dao"
	logger "github.com/RoyceAzure/sexy_gpt/broker_service/repository/logger_distributor"
	pawaiservicedaogo "github.com/RoyceAzure/sexy_gpt/broker_service/repository/pawai_service_dao"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/pb"
	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot start gateway http server")
	}
	setUpLoggerDistributor(config.RedisQueueAddress, config.ServiceID)

	accoutServiceDao, closeConnAccount, err := accountservicedao.NewAccountServiceDao(config.GrpcAccountAddress)
	if err != nil {
		logger.Logger.Fatal().
			Err(err).
			Msg("cannot connect account service grpc server")
	}
	defer closeConnAccount()

	pawAIServiceDao, closeConnPawAI, err := pawaiservicedaogo.NewPawAIServiceDao(config.GrpcPawAIAddress)
	if err != nil {
		logger.Logger.Fatal().
			Err(err).
			Msg("cannot connect pawai service grpc server")
	}
	defer closeConnPawAI()

	runGRPCGatewayServer(config, accoutServiceDao, pawAIServiceDao)

}

func setUpLoggerDistributor(address string, serviceId string) {
	redisOpt := asynq.RedisClientOpt{
		Addr: address,
	}

	//set up mongo logger
	redisClient := asynq.NewClient(redisOpt)
	loggerDis := logger.NewLoggerDistributor(redisClient)
	err := logger.SetUpLoggerDistributor(loggerDis, serviceId)
	if err != nil {
		log.Fatal().Err(err).Msg("err create mongo db connect")
	}
}

func runGRPCGatewayServer(
	configs config.Config,
	accoutServiceDao accountservicedao.IAccountServiceDao,
	pawAIServiceDao pawaiservicedaogo.IPawAIServiceDao,
) {

	tokenMakerAccount, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot start gateway http server")
	}

	tokenMakerPawAI, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot start gateway http server")
	}

	authorizeAccount := gapi.NewAuthorizor(tokenMakerAccount)

	authorizePaw := gapi.NewAuthorizor(tokenMakerPawAI)

	// 創建新的gRPC伺服器

	serverAccount := gapi.NewAccountServer(authorizeAccount, accoutServiceDao)

	serverPawAI := gapi.NewPawAIServerServer(authorizePaw, pawAIServiceDao, accoutServiceDao)

	jsonOpt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: false,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// 初始化gRPC Gateway的多路復用器
	/*runtime.NewServeMux() 創建的是一個 gRPC-Gateway 的多路復用器（multiplexer），它允許你將 HTTP/JSON 請求轉換為 gRPC 請求。

	它是一個 handler 嗎？

	是的，它實現了 http.Handler 接口，因此你可以將其用作 HTTP 伺服器的主要處理器。
	它是一個 multiplexer 嗎？

	是的，它是一個特殊的多路復用器，專為將 HTTP 請求轉換為 gRPC 請求而設計。當一個 HTTP 請求到達時，這個多路復用器會根據註冊的 gRPC 路由和方法轉換該請求，然後轉發它到對應的 gRPC 伺服器方法。
	總之，runtime.NewServeMux() 既是一個 handler，也是一個 multiplexer，但它專為 grpc-gateway 設計，用於在 gRPC 伺服器和 HTTP 客戶端之間進行轉換和路由。*/
	grpcMux := runtime.NewServeMux(jsonOpt, runtime.WithMetadata(middleware.CustomMatcher))

	// 創建一個可取消的背景上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 註冊gRPC伺服器到Gateway的多路復用器
	/*

		你在service.proto裡面定義的路由跟function, 都會在RegisterStockInfoHandlerServer 被設置，
		當呼叫RegisterStockInfoHandlerServer時，就會把路由以及handler設定到 *runtime.ServeMux上面
		RegisterStockInfoHandlerServer 會直接call grpc function (由RegisterStockInfoHandlerServer設置) , 不會經過intercepter
		RegisterStockInfoHandlerServer會把路由根handler 設置在你傳入的grpcMux 參數
	*/
	err = pb.RegisterAccountServiceHandlerServer(ctx, grpcMux, serverAccount)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot register account handler server")
	}

	err = pb.RegisterPawAIServiceHandlerServer(ctx, grpcMux, serverPawAI)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot register scheduler handler server")
	}

	// 創建HTTP多路復用器並將gRPC多路復用器掛載到其上
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	/*
			如果路由結尾有斜線 (/)，它將匹配任何以該前綴開始的 URL。因此，/swagger/ 將匹配 /swagger/、/swagger/file1.html、/swagger/subdir/file2.html 等。
		如果路由沒有結尾斜線，它只會匹配該具體路徑。
		所以，如果你使用 http.Handle("swagger/")（注意，缺少前置的斜線）：

		它將不會如預期地工作，因為在 net/http 中，路由通常應該以斜線 (/) 開始。這可能會導致未定義的行為或不匹配任何路徑。
		正確的做法是：

		使用 http.Handle("/swagger/") 以匹配所有以 /swagger/ 開始的路徑。
		使用 http.Handle("/swagger")（沒有結尾的斜線）只匹配 /swagger 這一具體的路徑，不匹配 /swagger/abc 或其他子路徑。
		總之，要確保路由以斜線 (/) 開始，並根據你的需求決定是否在結尾添加斜線。*/

	// fs := http.FileServer(http.Dir("./doc/swagger"))

	//這個FileSystem內容是zip content data, 剛好跟statik使用的依樣

	//在statik.go init()裡面已經註冊了使用statik編譯好的data
	//這裡New就是把他轉成fileSystem, 然後再搭配http.FileServer 即可
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("can't create statik fs err")
	}

	//http.StripPrefix 會回傳handler
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	// 在指定地址上建立監聽
	listener, err := net.Listen("tcp", configs.HttpServerAddress)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot create listener")
	}
	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())

	//
	loggerHandler := middleware.HttpLogger(mux)
	handler := middleware.IdMiddleWareHandler(loggerHandler)
	// handler1 := gapi.IdMiddleWareHandler(mux)
	// handler := gapi.HttpLogger(handler1)

	// 啟動HTTP伺服器
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot start HTTP gateway server")
	}
}
