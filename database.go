package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

type category struct {
	id          int
	description string
}

type budget struct {
	client_id int
	cat_id    int
	amount    sql.NullFloat64
}

// initializes the database. panics if a failure.
func init() {
	var err error
	// do NOT use in production
	connStr := "user=postgres password=? dbname=ecocosts sslmode=disable"
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}
