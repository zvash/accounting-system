package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zvash/accounting-system/internal/util"
	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
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
		return
	}
	fmt.Println(values)

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
