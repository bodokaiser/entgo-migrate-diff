package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/bodokaiser/entgo-migrate-diff/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var url string
var name string

func main() {
	flag.StringVar(&url, "url", "", "url of the postgres database")
	flag.StringVar(&name, "name", "", "name of the diff")
	flag.Parse()

	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("failed opening database connection: %v", err)
	}

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	defer client.Close()

	dir, err := sqltool.NewGolangMigrateDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating migration directory: %v", err)
	}

	err = client.Schema.NamedDiff(context.Background(), name,
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.Postgres),
	)
	if err != nil {
		log.Fatalf("failed creating named schema diff: %v", err)
	}
}
