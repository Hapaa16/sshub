package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func NewSQlclient() {
	var err error
	DBInstance, err := sql.Open("sqlite3", "./app.db")

	if err != nil {
		log.Fatal(err)
	}
	DBInstance.SetMaxOpenConns(1) // SQLite should have only 1 open connection

	DB = DBInstance

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
 CREATE TABLE IF NOT EXISTS connections (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  Title varchar(255),
  Host TEXT,
  Username varchar(255),
  IsPasswordConnection BOOL,
  Password TEXT,
  PemLocation TEXT
 );`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt) // Log an error if table creation fails
	}

}
