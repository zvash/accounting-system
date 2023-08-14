package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zvash/accounting-system/internal/api"
	"github.com/zvash/accounting-system/internal/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	serverAddress = "0.0.0.0:8080"
	dbSource      = "postgres://root:123@127.0.0.1:5432/simple_bank?sslmode=disable&pool_max_conns=32&pool_min_conns=10"
)

func main() {
	server := api.NewServer(createDBConnectionPool())
	log.Fatal(server.Start(serverAddress))
}

func createDBConnectionPool() *sql.DBStore {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	res, err := connPool.Query(context.Background(), "SHOW TRANSACTION ISOLATION LEVEL;")
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}
	res.Next()
	values, err := res.Values()
	if err != nil {
		return nil
	}
	fmt.Println(values)
	return sql.NewStore(connPool)
}
