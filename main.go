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
	server := api.NewServer(createDBConnectionPool(config.DBSource))
	log.Fatal(server.Start(config.HTTPServerAddress))
}

func createDBConnectionPool(dbSource string) *sql.DBStore {
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
