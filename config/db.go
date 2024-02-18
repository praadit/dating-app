package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func DB(conf Configuration) *bun.DB {
	sqlxDB, err := sqlx.Open("postgres", conf.DBConn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	sqlxDB.SetConnMaxLifetime(2 * time.Minute)
	Migrate(sqlxDB.DB)

	bunDB := bun.NewDB(sqlxDB.DB, pgdialect.New())
	// if config.Env == "development" {
	bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	// }

	return bunDB
}

func Migrate(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	migrate.SetTable("../migrations")

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Printf("Error applying db migration!\n%s", err)
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations!\n", n)
	}
}
