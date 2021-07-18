package models

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Conn struct {
	db *sqlx.DB
}

func Connect() Conn {
	dbPassword := os.Getenv("PSQL_PASS")

	if dbPassword == "" {
		log.Fatal("Database password not set")
	}

	connString := fmt.Sprintf("user=chidi password=%v dbname=chapi sslmode=disable", dbPassword)

	db, err := sqlx.Connect("postgres", connString)

	if err != nil {
		log.Fatalf("couldn't establish db connection: %v", err)
	}

	db.MustExec(schema)

	return Conn{db: db}
}
