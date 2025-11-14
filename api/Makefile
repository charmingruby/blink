.PHONY: protogen
protogen:
	find ./proto -name '*.proto' -print0 | xargs -0 protoc \
		--proto_path=./proto \
		--go_out=./proto/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto/pb \
		--go-grpc_opt=paths=source_relative
