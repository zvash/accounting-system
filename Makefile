postgres:
	docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -p 5432:5432 -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/sql/migration -database "postgresql://root:123@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

mr:
	make migratedown && make migrateup

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown mr sqlc test
