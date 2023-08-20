postgres:
	docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -p 5432:5432 -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

redis:
	docker run --name redis -p 6379:6379 -d redis:7.0-alpine

mu:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

mu1:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up 1

md:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

md1:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down 1

mr:
	make md && make mu

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb --destination internal/sql/mock/store.go github.com/zvash/accounting-system/internal/sql Store

proto:
	rm internal/pb/*
	protoc --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative --proto_path=internal/proto internal/proto/*.proto

.PHONY: postgres createdb dropdb mu md mr mu1 md1 sqlc test server mock proto redis
