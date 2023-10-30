package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	*TodoRepository
}

func NewPostgresRepository() *Repository {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")
	// postgres://<dbUsername>:<dbPassword>@<dbHost>/<dbName>?sslmode=<dbSslMode>
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbName, dbSslMode)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error", err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Connected to db!")

	return &Repository{
		NewTodoRepository(db),
	}
}
