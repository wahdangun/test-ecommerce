package database

import (
	"fmt"
	"os"

	"github.com/create-go-app/fiber-go-template/app/queries"
	"github.com/jmoiron/sqlx"
)

// Queries struct for collect all app queries.
type Queries struct {
	*queries.UserQueries    // load queries from User model
	*queries.BookQueries    // load queries from Book model
	*queries.ProductQueries // load queries from Product model
	*queries.CartQueries    // load queries from Cart model
	*queries.InvoiceQueries // load queries from Invoice model
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define Database connection variables.
	var (
		db  *sqlx.DB
		err error
	)

	// Get DB_TYPE value from .env file.
	dbType := os.Getenv("DB_TYPE")

	// Define a new Database connection with right DB type.
	switch dbType {
	case "pgx":
		db, err = PostgreSQLConnection()
	case "mysql":
		db, err = MysqlConnection()
	}

	if err != nil {
		return nil, err
	}
	fmt.Println("alamat", db)
	return &Queries{
		// Set queries from models:
		UserQueries:    &queries.UserQueries{DB: db},    // from User model
		BookQueries:    &queries.BookQueries{DB: db},    // from Book model
		ProductQueries: &queries.ProductQueries{DB: db}, // from Product model
		CartQueries:    &queries.CartQueries{DB: db},    // from Cart model
		InvoiceQueries: &queries.InvoiceQueries{DB: db}, // from Invoice model
	}, nil
}
