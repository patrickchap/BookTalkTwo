package database

import (
	db "BookTalkTwo/db/sqlc"
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	store db.Store
}

type service struct {
	store db.Store
}

var (
	dburl = os.Getenv("DB_URL")
	ddl   string
)

func New() Service {
	ctx := context.Background()
	sqlDb, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	if _, err := sqlDb.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(sqlDb)
	s := Service{store: store}
	return s
}
