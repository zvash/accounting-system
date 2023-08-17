package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zvash/accounting-system/internal/api"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	store := createDBConnectionPool(config.DBSource)
	createHttpServer(config, store)
}

func createHttpServer(config util.Config, store sql.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}
	log.Fatal(server.Start(config.HTTPServerAddress))
}

func createDBConnectionPool(dbSource string) sql.Store {
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
