protoc --proto_path=. --go_out=:. ./*.proto --go-grpc_out=require_unimplemented_servers=false:.