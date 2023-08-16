postgres:
	docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -p 5432:5432 -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

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

.PHONY: postgres createdb dropdb mu md mr mu1 md1 sqlc test server mock
