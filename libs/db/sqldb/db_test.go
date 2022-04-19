package sqldb

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {

	config := Configuration{
		Key: "pgsql1",
		Dsn: "postgres://postgres:12345678@127.0.0.1:5432/postgres?sslmode=disable",
	}

	// Open once on start application
	err := Open(&config)
	if err != nil {
		t.Error(err)
		return
	}

	db := DB(config.Key)
	var ver sql.NullString
	if err := db.QueryRow("SELECT version()").Scan(&ver); err != nil {
		t.Error(err)
	}

	if ver.Valid {
		fmt.Printf("\t== %s\n", ver.String)
	}
	// Ensure
	// Close once before application shutdown
	Close()
}
