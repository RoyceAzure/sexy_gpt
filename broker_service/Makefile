current_dir = $(shell echo %cd%)
DB_URL_LOCAL=postgres://royce:royce@localhost:5432/sexy_gpt?sslmode=disable
# DB_URL_AWS=postgres://royce:gqD2yhIOpUpuwK6IX6xz@stockinfo.cblayv8xneas.ap-northeast-1.rds.amazonaws.com:5432/stockinfo

postgresup:
	docker run --name project7-postgres --network stockinfo-network -e POSTGRES_PASSWORD=royce -e POSTGRES_USER=royce -e POSTGRES_DB=stock_info -p 5432:5432 -v D:\Workspace\JackRabbit\GO\project7\project\db-data\postgres:/var/lib/postgresql/data/ -d postgres:14.2

postgresrm:
	docker stop project7-postgres
	docker rm project7-postgres
	rm -Recurse -Force .\db-data\postgres

createdb:
	docker exec -it project7-postgres  createdb --username=royce --owner=royce stock_info

dropdb:
	docker exec -it project7-postgres  dropdb --username=royce stock_info

sqlc:
	docker run --rm -v $(current_dir)/repository/db:/src -w /src sqlc/sqlc:latest generate

awsmigrateup:
	migrate -path repository/db/migrations -database $(DB_URL_AWS) --verbose up

migrateup:
	migrate -path repository/db/migrations -database $(DB_URL_LOCAL) --verbose up

migratedown:
	migrate -path repository/db/migrations -database $(DB_URL_LOCAL) --verbose down

migrateforce:
	migrate -path repository/db/migrations -database $(DB_URL_LOCAL) --verbose force $(stage)

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.2-alpine

test:
	go test -v -cover -short ./...

server:
	go run main.go

# mock:
# 	mockgen -source="./distributor.go" -package=mockwk -destination="mock/distributor_mock.go"
# 	mockgen  . Store
#mockgen --destination github.com/RoyceAzure/sexy_gpt/account_service/repository/db/mock --package mock_sqlc  github.com/RoyceAzure/sexy_gpt/account_service/repository/db/sqlc Dao
#mockgen --destination github.com/RoyceAzure/sexy_gpt/account_service/service/mock --package mock_service  github.com/RoyceAzure/sexy_gpt/account_service/service IService
#mockgen --destination github.com/RoyceAzure/sexy_gpt/account_service/worker/mock --package mock_worker  github.com/RoyceAzure/sexy_gpt/account_service/worker TaskDistributor
db_docs:
	docker run --rm -v $(current_dir)/doc:/app/data -w /app/data node_docs dbdocs build ./db.dbml

db_schema:
	docker run --rm -v $(current_dir)/doc:/app -w /app node_docs dbml2sql --postgres -o sexy_gpt_schema.sql db.dbml

protoc:
	powershell -Command "Remove-Item -Path 'shared/pb/*.go' -Force"
	protoc   --grpc-gateway_out shared/pb \
	-I=./proto -I=C:/Users/royce/include \
	--proto_path=proto  --go_out=shared/pb  --go_opt=paths=source_relative  --grpc-gateway_opt=paths=source_relative \
	--go-grpc_out=shared/pb --go-grpc_opt=paths=source_relative \
	--openapiv2_out doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=sexy_gpt \
	proto/*.proto

evans:
	docker run -it --rm -v $(current_dir):/mount:ro ghcr.io/ktr0731/evans:latest --host host.docker.internal --port 9090 -r repl

.PHONY: postgresup postgresrm createdb dropdb test server awsmigrateup db_docs db_schema protoc evans redis new_migration migrateforce migrateupto migratedownto
 