package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/services"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var Service *services.Service
var BunDB *bun.DB

func TestMain(m *testing.M) {
	BunDB = testDbSetup()
	memCache := cache.New(1*time.Minute, 10*time.Minute)
	Service = services.NewService(BunDB, memCache)

	os.Exit(m.Run())
}

func testDbSetup() *bun.DB {
	dbConn, found := os.LookupEnv("DB_CONN")
	if !found {
		err := godotenv.Load("../.env")
		if err != nil {
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}
		dbConn = os.Getenv("DB_CONN")
	} else {
		log.Println("using system envi variable")
	}

	fmt.Println("DB CONN", dbConn)
	sqlxDB, err := sqlx.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	sqlxDB.SetMaxOpenConns(4)
	sqlxDB.SetMaxIdleConns(4)
	sqlxDB.SetConnMaxLifetime(2 * time.Minute)

	return bun.NewDB(sqlxDB.DB, pgdialect.New())
}
