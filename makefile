grpc:
	go run users_grpc_server/main.go

api:
	go run users_api_server/main.go

generate:
	cd users_proto && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=. users.proto
